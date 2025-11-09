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

	db "github.com/ipincamp/srikandi-sehat/internal/adapters/driven/postgres"
)

// Compile-time check
var _ ports.AuthService = (*authService)(nil)

// authService implements the ports.AuthService interface.
type authService struct {
	userRepo  ports.UserRepository // For non-transactional reads (e.g., Login)
	userCache ports.UserCache
	maker     token.Maker
	hasher    password.Hasher
	tokenCfg  config.Token
	logger    zerolog.Logger
	uow       ports.UnitOfWork
}

// NewAuthService is the constructor for authService.
func NewAuthService(
	userRepo ports.UserRepository, // This is the non-transactional repo
	userCache ports.UserCache,
	maker token.Maker,
	hasher password.Hasher,
	tokenCfg config.Token,
	logger zerolog.Logger,
	uow ports.UnitOfWork,
) ports.AuthService {
	return &authService{
		userRepo:  userRepo,
		userCache: userCache,
		maker:     maker,
		hasher:    hasher,
		tokenCfg:  tokenCfg,
		logger:    logger,
		uow:       uow,
	}
}

// Register creates a new user, hashes their password,
// saves them, and returns a new set of auth tokens.
// This operation is now transactional.
func (s *authService) Register(ctx context.Context, name, email, passwordStr string) (*ports.AuthResponse, error) {
	// 1. Check if user *might* exist using the bloom filter
	if s.userCache.Test(email) {
		// Email *might* exist. We must fallback to the DB for a definitive check.
		// We use the non-transactional repo for this read-only check.
		_, err := s.userRepo.FindByEmail(ctx, email)
		if err == nil {
			// User found, email is taken
			s.logger.Warn().Str("email", email).Msg("Registration failed: email already exists (pre-check)")
			return nil, ErrEmailExists
		}
		if !errors.Is(err, db.ErrUserNotFound) {
			// A different, unexpected database error occurred during find
			s.logger.Error().Err(err).Str("email", email).Msg("Failed to check user existence")
			return nil, err
		}
		// If we are here, it was a false positive. We can proceed.
	}
	// If filter.Test() was false, email *definitely does not exist*,
	// so we skip the FindByEmail check entirely, saving a DB query.

	// 2. Hash the password
	hashedPassword, err := s.hasher.Hash(passwordStr)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to hash password during registration")
		return nil, err
	}

	// 3. Create the domain user
	user := &domain.User{
		UUID:     uuid.NewString(),
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		// ID, CreatedAt, UpdatedAt will be set by the repository
	}

	// 4. === Begin Transactional Unit of Work ===
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

	// 5. Get the transactional repository from the Unit of Work
	txUserRepo := tx.GetUserRepository()

	// 6. Save the user *using the transactional repo*
	if err = txUserRepo.Save(ctx, user); err != nil {
		// Check for duplicate email (race condition)
		if errors.Is(err, db.ErrDuplicateEmail) {
			s.logger.Warn().Str("email", email).Msg("Registration failed: email already exists (race condition on save)")
			return nil, ErrEmailExists // Defer will catch this and rollback
		}

		// A different, unexpected save error
		s.logger.Error().Err(err).Str("email", email).Msg("Failed to save user during registration")
		return nil, err // Defer will catch this and rollback
	}

	// 7. (Example) If you had other tables, you would save them here
	// e.g., tx.GetProfileRepository().CreateDefaultProfile(ctx, user.ID)
	// If this failed, the defer would roll back the user creation.

	// 8. Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Failed to commit registration transaction")
		return nil, err // err is already set, so defer will *not* roll back again
	}
	// === End Transactional Unit of Work ===

	s.logger.Info().Str("email", email).Str("uuid", user.UUID).Msg("User registered successfully")

	// 5. Add new user to our in-memory cache
	s.userCache.Add(user.Email)

	// 6. Generate tokens
	return s.createTokenSet(user)
}

// Login validates user credentials and returns a new set of auth tokens.
func (s *authService) Login(ctx context.Context, email, passwordStr string) (*ports.AuthResponse, error) {
	// 1. Check bloom filter
	if !s.userCache.Test(email) {
		// Email *definitely does not exist*.
		// We can short-circuit without hitting the DB.
		// This is a very cheap way to reject invalid login attempts.
		s.logger.Warn().Str("email", email).Msg("Login failed: invalid credentials (user not found via bloom filter)")
		return nil, ErrInvalidCredentials
	}

	// 2. Find user by email (filter reported a *possible* match)
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			s.logger.Warn().Str("email", email).Msg("Login failed: invalid credentials (user not found)")
			return nil, ErrInvalidCredentials
		}
		s.logger.Error().Err(err).Str("email", email).Msg("Login failed: database error on find")
		return nil, err
	}

	// 3. Compare password
	if !s.hasher.Compare(user.Password, passwordStr) {
		s.logger.Warn().Str("email", email).Msg("Login failed: invalid credentials (password mismatch)")
		return nil, ErrInvalidCredentials
	}

	// 4. Generate tokens
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
			return nil, ErrTokenExpired
		}
		s.logger.Warn().Err(err).Msg("Refresh token failed: invalid token")
		return nil, ErrInvalidToken
	}

	// 2. Check that it's actually a refresh token
	if payload.UseFor != token.UseForRefreshToken {
		s.logger.Warn().Str("uuid", payload.UserID).Msg("Refresh token failed: token use mismatch")
		return nil, ErrTokenUseMismatch
	}

	// 3. Find the user
	user, err := s.userRepo.FindByID(ctx, payload.UserID)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			s.logger.Error().Str("uuid", payload.UserID).Msg("Refresh token failed: user not found")
			return nil, ErrUserNotFound
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
