package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/ipincamp/srikandi-sehat/pkg/config"
	"github.com/ipincamp/srikandi-sehat/pkg/password"
	"github.com/ipincamp/srikandi-sehat/pkg/token"
)

// Compile-time check
var _ ports.AuthService = (*authService)(nil)

// authService implements the ports.AuthService interface.
type authService struct {
	userRepo ports.UserRepository // For non-transactional reads (e.g., Login)
	maker    token.Maker
	hasher   password.Hasher
	tokenCfg config.Token
	logger   zerolog.Logger
	uow      ports.UnitOfWork

	// --- Dependensi Blueprint ---
	// Kita suntikkan interface-nya, meskipun implementasinya belum ada.
	// emailSvc ports.EmailServicePort
	// otpSvc   ports.OTPServicePort
}

// NewAuthService is the constructor for authService.
func NewAuthService(
	userRepo ports.UserRepository, // This is the non-transactional repo
	maker token.Maker,
	hasher password.Hasher,
	tokenCfg config.Token,
	logger zerolog.Logger,
	uow ports.UnitOfWork,
) ports.AuthService {
	return &authService{
		userRepo: userRepo,
		maker:    maker,
		hasher:   hasher,
		tokenCfg: tokenCfg,
		logger:   logger,
		uow:      uow,
	}
}

// Register creates a new user, hashes their password,
// saves them, and returns a new set of auth tokens.
// This operation is now transactional.
func (s *authService) Register(ctx context.Context, name, email, passwordStr string) (*ports.AuthResponse, error) {
	/*
		DEPRECATED
		FindByEmail check outside of transactions
		We will rely on unique database constraints within transactions
		to handle email duplication (eliminating race conditions).
	*/

	// 1. Hash the password
	hashedPassword, err := s.hasher.Hash(passwordStr)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to hash password during registration")
		return nil, err
	}

	// 2. Create the domain user
	user := &domain.User{
		UUID:     uuid.NewString(),
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		// ID, CreatedAt, UpdatedAt will be set by the repository
	}

	// 3. === Begin Transactional Unit of Work ===
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to begin registration transaction")
		return nil, err
	}

	// Defer a function to handle rollback in case of panic or error
	defer func() {
		if r := recover(); r != nil {
			// A panic occurred
			s.logger.Error().Msgf("Panic detected in Register, rolling back transaction: %v", r)
			_ = tx.Rollback(ctx)
			panic(r) // re-panic after rollback
		}
		if err != nil {
			// An error occurred, rollback the transaction
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				s.logger.Error().Err(rbErr).Msg("Failed to rollback transaction after error")
			}
		}
	}()

	// 4. Get the transactional repository from the Unit of Work
	txUserRepo := tx.GetUserRepository()

	// 5. Save the user *using the transactional repo*
	if err = txUserRepo.Save(ctx, user); err != nil {
		// Check for duplicate email (race condition)
		if errors.Is(err, ports.ErrDuplicateEmail) {
			s.logger.Warn().Str("email", email).Msg("Registration failed: email already exists (constraint violation)")
			return nil, ports.ErrEmailExists // Defer will catch this and rollback
		}

		// A different, unexpected save error
		s.logger.Error().Err(err).Str("email", email).Msg("Failed to save user during registration")
		return nil, err // Defer will catch this and rollback
	}

	// 6. (Example) If you had other tables, you would save them here
	// e.g., tx.GetProfileRepository().CreateDefaultProfile(ctx, user.ID)
	// If this failed, the defer would roll back the user creation.

	// 7. Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Failed to commit registration transaction")
		return nil, err // err is already set, so defer will *not* roll back again
	}
	// === End Transactional Unit of Work ===

	s.logger.Info().Str("email", email).Str("uuid", user.UUID).Msg("User registered successfully")

	// 8. Generate tokens
	return s.createTokenSet(user)
}

// Login validates user credentials and returns a new set of auth tokens.
func (s *authService) Login(ctx context.Context, email, passwordStr string) (*ports.AuthResponse, error) {
	// 1. Find user by email (we now *always* hit the DB for this).
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			s.logger.Warn().Msg("Login failed: invalid credentials (user not found)")
			return nil, ports.ErrInvalidCredentials
		}
		s.logger.Error().Err(err).Str("email", email).Msg("Login failed: database error on find")
		return nil, err
	}

	// 2. Compare password
	if !s.hasher.Compare(user.Password, passwordStr) {
		s.logger.Warn().Msg("Login failed: invalid credentials (password mismatch)")
		return nil, ports.ErrInvalidCredentials
	}

	// 3. Generate tokens
	s.logger.Info().Str("email", email).Str("uuid", user.UUID).Msg("User logged in successfully")
	return s.createTokenSet(user)
}

// RefreshToken validates a refresh token and issues a new pair of tokens.
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*ports.AuthResponse, error) {
	// 1. Validate the refresh token
	payload, err := s.maker.ValidateToken(refreshToken)
	if err != nil {
		if errors.Is(err, token.ErrTokenExpired) {
			s.logger.Warn().Msg("Refresh token failed: token expired")
			return nil, ports.ErrTokenExpired
		}
		s.logger.Warn().Err(err).Msg("Refresh token failed: invalid token")
		return nil, ports.ErrInvalidToken
	}

	// 2. Check that it's actually a refresh token
	if payload.UseFor != token.UseForRefreshToken {
		s.logger.Warn().Str("uuid", payload.UserID).Msg("Refresh token failed: token use mismatch")
		return nil, ports.ErrTokenUseMismatch
	}

	// 3. Find the user
	user, err := s.userRepo.FindByID(ctx, payload.UserID)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			s.logger.Error().Str("uuid", payload.UserID).Msg("Refresh token failed: user not found")
			return nil, ports.ErrUserNotFound
		}
		s.logger.Error().Err(err).Str("uuid", payload.UserID).Msg("Refresh token failed: database error")
		return nil, err
	}

	// 4. Generate new tokens
	s.logger.Info().Str("uuid", user.UUID).Msg("Token refreshed successfully")
	return s.createTokenSet(user)
}

// Logout is a no-op for stateless tokens.
// The client is responsible for deleting the tokens.
// If we had a stateful repository (e.g., in Redis or DB),
// we would invalidate the refresh token here.
func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	// For a stateful implementation:
	// 1. Validate token
	// 2. Get payload
	// 3. Call `refreshTokenRepo.Delete(ctx, payload.TokenID)`
	// 4. Return result

	// For this stateless implementation:
	return nil
}

func (s *authService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	// 1. Dapatkan user (non-transaksional untuk pengecekan)
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err // Termasuk ErrUserNotFound
	}

	// 2. Verifikasi password lama
	if !s.hasher.Compare(user.Password, oldPassword) {
		return ports.ErrInvalidCredentials
	}

	// 3. Hash password baru
	newHashedPassword, err := s.hasher.Hash(newPassword)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to hash new password")
		return err
	}

	// 4. Mulai Unit of Work untuk menyimpan
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to begin ChangePassword transaction")
		return err
	}
	defer tx.Rollback(ctx) // Rollback jika ada error

	// 5. Dapatkan repo transaksional dan update
	txUserRepo := tx.GetUserRepository()
	user.Password = newHashedPassword
	// Kita perlu memastikan FindByID di-implementasikan oleh repo transaksional
	// atau kita perlu mem-fetch ulang user di dalam transaksi.
	// Untuk saat ini, kita asumsikan Update bisa menangani ini.
	// (Cara lebih aman: fetch user *di dalam* transaksi)
	if err := txUserRepo.Update(ctx, user); err != nil {
		s.logger.Error().Err(err).Str("uuid", userID).Msg("Failed to update password in DB")
		return err
	}

	// 6. Commit
	if err := tx.Commit(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Failed to commit ChangePassword transaction")
		return err
	}

	s.logger.Info().Str("uuid", userID).Msg("Password changed successfully")
	return nil
}

// --- Implementasi Blueprint ---

// 5. ForgotPassword (Blueprint)
func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Jangan bocorkan apakah email ada atau tidak
		s.logger.Warn().Str("email", email).Msg("Forgot password attempt for non-existent email")
		return nil // Selalu return nil
	}

	// const OTPSDuration = 15 * time.Minute
	// otp, err := s.otpSvc.GenerateAndStoreOTP(ctx, user.UUID, "password_reset", OTPSDuration)
	// if err != nil {
	// 	s.logger.Error().Err(err).Msg("Failed to generate OTP for password reset")
	// 	return err
	// }
	//
	// err = s.emailSvc.SendPasswordResetEmail(ctx, user.Email, user.Name, otp)
	// if err != nil {
	// 	s.logger.Error().Err(err).Msg("Failed to send password reset email")
	// 	return err
	// }

	s.logger.Info().Str("uuid", user.UUID).Msg("Forgot password process initiated (blueprint)")
	return nil // TODO: Hapus blueprint stub
}

// 6. VerifyEmailOTP (Blueprint)
func (s *authService) VerifyEmailOTP(ctx context.Context, otp string) error {
	// userID, err := s.otpSvc.ValidateAndConsumeOTP(ctx, otp, "email_verification")
	// if err != nil {
	// 	s.logger.Warn().Err(err).Msg("Failed to validate email OTP")
	// 	return ports.ErrInvalidToken
	// }
	//
	// tx, err := s.uow.Begin(ctx)
	// ... (logika untuk fetch user by ID, set user.IsVerified = true, txUserRepo.Update(user), tx.Commit()) ...

	s.logger.Info().Str("otp", otp).Msg("Email verification attempt (blueprint)")
	return nil // TODO: Hapus blueprint stub
}

// 10. RequestEmailChange (Blueprint)
func (s *authService) RequestEmailChange(ctx context.Context, userID, newEmail string) error {
	// ... (logika cek jika email baru sudah dipakai) ...
	// ... (logika panggil s.otpSvc.GenerateAndStoreOTP) ...
	// ... (logika panggil s.emailSvc.SendEmailChangeEmail) ...
	s.logger.Info().Str("uuid", userID).Str("newEmail", newEmail).Msg("Email change requested (blueprint)")
	return nil // TODO: Hapus blueprint stub
}

// createTokenSet is a helper to generate both access and refresh tokens.
func (s *authService) createTokenSet(user *domain.User) (*ports.AuthResponse, error) {
	// Create Access Token
	accessToken, _, err := s.maker.CreateToken(
		user.UUID,
		"", // RoleID - not implemented yet
		token.UseForAccessToken,
		s.tokenCfg.AccessTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	// Create Refresh Token
	refreshToken, _, err := s.maker.CreateToken(
		user.UUID,
		"", // RoleID
		token.UseForRefreshToken,
		s.tokenCfg.RefreshTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	return &ports.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
