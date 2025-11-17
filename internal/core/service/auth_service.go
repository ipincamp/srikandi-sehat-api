package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/ipincamp/srikandi-sehat/pkg/config"
	"github.com/ipincamp/srikandi-sehat/pkg/password"
	"github.com/ipincamp/srikandi-sehat/pkg/token"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Ensure authService implements the port
var _ ports.AuthService = (*authService)(nil)

type authService struct {
	uow         ports.UnitOfWork
	hasher      password.Hasher
	tokenMaker  token.Maker
	tokenCfg    config.Token
	logger      zerolog.Logger
	mailService ports.MailService
}

// NewAuthService is the constructor for AuthService
func NewAuthService(
	uow ports.UnitOfWork,
	hasher password.Hasher,
	tokenMaker token.Maker,
	tokenCfg config.Token,
	logger zerolog.Logger,
	mailService ports.MailService,
) ports.AuthService {
	return &authService{
		uow:         uow,
		hasher:      hasher,
		tokenMaker:  tokenMaker,
		tokenCfg:    tokenCfg,
		logger:      logger,
		mailService: mailService,
	}
}

// handleRollback is a helper to safely rollback a transaction and log failures
func (s *authService) handleRollback(tx ports.Transaction, logMsg string) {
	if err := tx.Rollback(); err != nil {
		s.logger.Error().Err(err).Msg(logMsg)
	}
}

// generateSecureToken creates a URL-safe, random string.
func generateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Use URL-safe encoding to avoid issues with email clients.
	return base64.URLEncoding.EncodeToString(b), nil
}

// hashToken uses SHA256 to hash the one-time-use token.
// We don't use Argon2 (s.hasher) as it's too slow for this purpose.
// A fast hash is appropriate here because the token's entropy comes
// from its high randomness (32 bytes), not a user-generated password.
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Register implements the registration logic
func (s *authService) Register(ctx context.Context, name, email, password string) (*domain.AuthResponse, error) {
	// 1. Start Transaction
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to begin transaction for Register")
		return nil, errors.New("registration failed")
	}

	// Get transactional repositories
	userRepo := tx.GetUserRepository()

	// 2. Check if user already exists (within tx)
	_, err = userRepo.FindByEmail(ctx, email)
	if err == nil {
		s.handleRollback(tx, "Rollback Register: email already in use")
		return nil, errors.New("email already in use")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.handleRollback(tx, "Rollback Register: failed to check email")
		s.logger.Error().Err(err).Msg("Failed to check user by email")
		return nil, errors.New("registration failed")
	}
	// If err is gorm.ErrRecordNotFound, we can proceed

	// 3. Hash the password
	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		s.handleRollback(tx, "Rollback Register: failed to hash password")
		s.logger.Error().Err(err).Msg("Failed to hash password")
		return nil, errors.New("registration failed")
	}

	// 4. Create the domain user
	newUser := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	// 5. Save the user (within tx)
	if err := userRepo.Save(ctx, newUser); err != nil {
		s.handleRollback(tx, "Rollback Register: failed to save user")
		s.logger.Error().Err(err).Msg("Failed to save new user")
		return nil, errors.New("registration failed")
	}

	// 6. Fetch the newly created user to get the ID (within tx)
	createdUser, err := userRepo.FindByEmail(ctx, email)
	if err != nil {
		s.handleRollback(tx, "Rollback Register: failed to fetch new user")
		s.logger.Error().Err(err).Msg("Failed to fetch newly created user")
		return nil, errors.New("registration failed")
	}

	// 7. Generate tokens and save JTI (within tx)
	authResponse, err := s.createTokenSet(ctx, tx, createdUser)
	if err != nil {
		s.handleRollback(tx, "Rollback Register: failed to create token set")
		return nil, err // Error already wrapped in createTokenSet
	}

	// 8. Commit Transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error().Err(err).Msg("Failed to commit transaction for Register")
		return nil, errors.New("registration failed")
	}

	// 9. Send welcome email asynchronously
	// This happens *after* the transaction is committed.
	// We run this in a goroutine so it doesn't block the user's response.
	go func() {
		// We create a new background context for the goroutine.
		emailCtx := context.Background()
		subject := "Welcome to Srikandi Sehat!"
		plainBody := fmt.Sprintf("Hi %s,\n\nWelcome! Please verify your email. (OTP logic to be added).", name)
		htmlBody := fmt.Sprintf("<h1>Hi %s,</h1><p>Welcome! Please verify your email. (OTP logic to be added).</p>", name)

		// The service calls the interface, completely unaware of "SMTP".
		if err := s.mailService.Send(emailCtx, email, subject, plainBody, htmlBody); err != nil {
			// Log the error, but don't return it to the user who already registered.
			s.logger.Error().Err(err).Str("user_email", email).Msg("Failed to send welcome email")
		}
	}()

	return authResponse, nil
}

// Login implements the login logic
func (s *authService) Login(ctx context.Context, email, password string) (*domain.AuthResponse, error) {
	// 1. Start Transaction
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to begin transaction for Login")
		return nil, errors.New("login failed")
	}
	// Add logger with context for this request
	log := s.logger.With().Str("method", "Login").Str("email", email).Logger()

	// Get transactional repositories
	userRepo := tx.GetUserRepository()

	// 2. Find the user by email (within tx)
	user, err := userRepo.FindByEmail(ctx, email)
	if err != nil {
		s.handleRollback(tx, "Rollback Login: find by email failed")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		log.Error().Err(err).Msg("Failed to find user by email for login")
		return nil, errors.New("login failed")
	}

	// 3. Compare the password
	if !s.hasher.Compare(user.PasswordHash, password) {
		s.handleRollback(tx, "Rollback Login: invalid password")
		return nil, errors.New("invalid email or password")
	}

	// 4. Check if account is disabled (Requirement 1.2.3.6)
	if user.IsDisabled() { // This method is from your domain/user.go file
		s.handleRollback(tx, "Rollback Login: account is disabled")
		log.Warn().Msg("Login attempt from disabled account")
		// This specific error message maps to requirement 1.2.3.6
		return nil, errors.New("invalid email or password")
	}

	// 5. Generate tokens and save JTI (within tx)
	authResponse, err := s.createTokenSet(ctx, tx, user)
	if err != nil {
		s.handleRollback(tx, "Rollback Login: failed to create token set")
		return nil, err
	}

	// 6. Commit Transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction for Login")
		return nil, errors.New("login failed")
	}

	return authResponse, nil
}

// createTokenSet is a helper that requires a transaction
func (s *authService) createTokenSet(ctx context.Context, tx ports.Transaction, user *domain.User) (*domain.AuthResponse, error) {
	// Get transactional repository from the transaction
	personalTokenRepo := tx.GetPersonalTokenRepository()

	// Create Access Token
	accessToken, _, err := s.tokenMaker.CreateToken(
		user.ID,
		"", // RoleID
		token.UseForAccessToken,
		s.tokenCfg.AccessTokenTTL,
	)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create access token")
		return nil, errors.New("token generation failed")
	}

	// Create Refresh Token
	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.ID,
		"", // RoleID
		token.UseForRefreshToken,
		s.tokenCfg.RefreshTokenTTL,
	)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create refresh token")
		return nil, errors.New("token generation failed")
	}

	// Save the refresh token JTI (within tx)
	personalToken := &domain.PersonalToken{
		ID:        refreshPayload.JTI,
		UserID:    user.ID,
		ExpiresAt: refreshPayload.ExpiresAt,
	}

	if err := personalTokenRepo.Save(ctx, personalToken); err != nil {
		s.logger.Error().Err(err).Msg("Failed to save personal token JTI")
		return nil, fmt.Errorf("token generation failed: %w", err)
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken validates an old refresh token and issues a new pair.
func (s *authService) RefreshToken(ctx context.Context, tokenString string) (*domain.AuthResponse, error) {
	// 1. Validate token (this is outside a transaction, it's just a check)
	payload, err := s.tokenMaker.ValidateToken(tokenString)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	if payload.UseFor != token.UseForRefreshToken {
		return nil, errors.New("invalid token purpose")
	}

	// 2. Start Transaction
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to begin transaction for RefreshToken")
		return nil, errors.New("token refresh failed")
	}

	// Get transactional repositories
	personalTokenRepo := tx.GetPersonalTokenRepository()
	userRepo := tx.GetUserRepository()

	// 3. Check if JTI is in DB (within tx)
	if _, err := personalTokenRepo.FindByID(ctx, payload.JTI); err != nil {
		s.handleRollback(tx, "Rollback Refresh: JTI not found")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn().Str("jti", payload.JTI).Msg("Refresh token JTI not found in DB")
			return nil, errors.New("invalid refresh token")
		}
		s.logger.Error().Err(err).Msg("Failed to find personal token by ID")
		return nil, errors.New("token refresh failed")
	}

	// 4. JTI is valid. Revoke it (Token Rotation) (within tx)
	if err := personalTokenRepo.Delete(ctx, payload.JTI); err != nil {
		s.handleRollback(tx, "Rollback Refresh: failed to delete old JTI")
		s.logger.Error().Err(err).Msg("Failed to delete old JTI during refresh")
		return nil, errors.New("token refresh failed")
	}

	// 5. Fetch the user (within tx)
	user, err := userRepo.FindByID(ctx, payload.UserID)
	if err != nil {
		s.handleRollback(tx, "Rollback Refresh: user not found")
		s.logger.Error().Err(err).Str("user_id", payload.UserID).Msg("User for refresh token not found")
		return nil, errors.New("invalid refresh token")
	}

	// 6. Create new token set (within tx)
	authResponse, err := s.createTokenSet(ctx, tx, user)
	if err != nil {
		s.handleRollback(tx, "Rollback Refresh: failed to create new token set")
		return nil, err
	}

	// 7. Commit Transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error().Err(err).Msg("Failed to commit transaction for RefreshToken")
		return nil, errors.New("token refresh failed")
	}

	return authResponse, nil
}

// Logout revokes a specific refresh token.
func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	// 1. Validate token (outside tx)
	payload, err := s.tokenMaker.ValidateToken(refreshToken)
	if err != nil {
		// If token is already invalid/expired, it's effectively "logged out"
		s.logger.Warn().Err(err).Msg("Logout attempt with invalid token")
		return nil
	}
	if payload.UseFor != token.UseForRefreshToken {
		s.logger.Warn().Msg("Logout attempt with non-refresh token")
		return nil // Client is confused, but no error needed
	}

	// 2. Start Transaction
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to begin transaction for Logout")
		return errors.New("logout failed") // This is a server error
	}

	// Get transactional repository
	personalTokenRepo := tx.GetPersonalTokenRepository()

	// 3. Delete JTI (within tx)
	// We must check if it exists first
	if _, err := personalTokenRepo.FindByID(ctx, payload.JTI); err == nil {
		// Found it, now delete it
		if err := personalTokenRepo.Delete(ctx, payload.JTI); err != nil {
			s.handleRollback(tx, "Rollback Logout: failed to delete JTI")
			s.logger.Error().Err(err).Str("jti", payload.JTI).Msg("Failed to delete JTI on logout")
			return errors.New("logout failed")
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// A database error occurred trying to find the token
		s.handleRollback(tx, "Rollback Logout: failed to find JTI")
		s.logger.Error().Err(err).Str("jti", payload.JTI).Msg("Failed to find JTI on logout")
		return errors.New("logout failed")
	}
	// If ErrRecordNotFound, the token is already revoked, which is fine.

	// 4. Commit Transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error().Err(err).Msg("Failed to commit transaction for Logout")
		return errors.New("logout failed")
	}

	return nil
}

// ForgotPassword implements the logic from section 1.5.3.
func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	// This logic *always* returns a nil error to prevent email enumeration.
	// The actual work (DB operations, email) happens only if the user is found.
	log := s.logger.With().Str("method", "ForgotPassword").Str("email", email).Logger()

	// We must use a transaction to get repositories.
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return nil // Still return nil to user
	}

	userRepo := tx.GetUserRepository()
	user, err := userRepo.FindByEmail(ctx, email)

	// If user not found  or any other DB error, just commit and return.
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("Failed to query user by email")
		}
		// Commit the (empty) transaction and return.
		_ = tx.Commit()
		return nil
	}

	// --- User was found, proceed with logic  ---
	tokenRepo := tx.GetUserTokenRepository()

	// 1. Generate new token.
	// We create a 32-byte random token.
	tokenString, err := generateSecureToken(32)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate secure token")
		s.handleRollback(tx, "Rollback ForgotPassword: failed to generate token")
		return nil
	}

	// 2. Hash the token for storage.
	tokenHash := hashToken(tokenString)
	tokenExpiry := time.Now().Add(15 * time.Minute) // 15-minute expiry.

	// 3. Save the new token hash to the DB.
	// This now uses our Upsert logic.
	userToken := &domain.UserToken{
		UserID:    user.ID,
		Purpose:   domain.TokenPurposePasswordReset,
		TokenHash: tokenHash,
		ExpiresAt: tokenExpiry,
	}
	if err := tokenRepo.Save(ctx, userToken); err != nil {
		log.Error().Err(err).Msg("Failed to save new reset token")
		s.handleRollback(tx, "Rollback ForgotPassword: failed to upsert token")
		return nil
	}

	// 4. Commit the transaction.
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return nil
	}

	// 5. Send email asynchronously.
	go func() {
		emailCtx := context.Background()
		log.Info().Msg("Dispatching password reset email")

		// The token sent to the user is in the format `uid.tokenstring`.
		fullToken := fmt.Sprintf("%s.%s", user.ID, tokenString)

		subject := "Your Password Reset Instructions"
		plainBody := fmt.Sprintf("Hi %s,\n\nYou requested a password reset. Use this token (it will expire in 15 minutes):\n\n%s\n\nIf you did not request this, please ignore this email.", user.Name, fullToken)
		htmlBody := fmt.Sprintf("<h1>Hi %s,</h1><p>You requested a password reset. Use this token (it will expire in 15 minutes):</p><h2>%s</h2><p>If you did not request this, please ignore this email.</p>", user.Name, fullToken)

		if err := s.mailService.Send(emailCtx, user.Email, subject, plainBody, htmlBody); err != nil {
			log.Error().Err(err).Msg("Failed to send password reset email")
		}
	}()

	return nil
}

// ResetPassword implements the logic from section 1.6.3.
func (s *authService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	log := s.logger.With().Str("method", "ResetPassword").Logger()

	// 1. Validate token format: "uid.tokenstring".
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return errors.New("invalid token format")
	}
	userID, tokenString := parts[0], parts[1]

	if userID == "" || tokenString == "" {
		return errors.New("invalid token format")
	}

	// 2. Start Transaction.
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return errors.New("password reset failed")
	}

	// 3. Find token in DB.
	tokenRepo := tx.GetUserTokenRepository()
	userToken, err := tokenRepo.FindByUserIDAndPurpose(ctx, userID, domain.TokenPurposePasswordReset)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msg("Password reset token not found in DB")
		} else {
			log.Error().Err(err).Msg("Failed to find token")
		}
		s.handleRollback(tx, "Rollback ResetPassword: token not found")
		return errors.New("invalid or expired token")
	}

	// 4. Hash the token from the user.
	tokenHash := hashToken(tokenString)

	// 5. Validate the token hash and expiry.
	// Use constant-time compare to prevent timing attacks.
	if subtle.ConstantTimeCompare([]byte(userToken.TokenHash), []byte(tokenHash)) != 1 {
		log.Warn().Msg("Token hash mismatch")
		s.handleRollback(tx, "Rollback ResetPassword: token hash mismatch")
		return errors.New("invalid or expired token")
	}

	// Check expiry.
	if time.Now().After(userToken.ExpiresAt) {
		log.Warn().Msg("Token expired")
		// Delete the expired token.
		_ = tokenRepo.DeleteByUserIDAndPurpose(ctx, userID, domain.TokenPurposePasswordReset)
		s.handleRollback(tx, "Rollback ResetPassword: token expired")
		return errors.New("invalid or expired token")
	}

	// --- Token is valid ---

	// 6. Update user's password.
	// Hash the new password.
	hashedPassword, err := s.hasher.Hash(newPassword)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash new password")
		s.handleRollback(tx, "Rollback ResetPassword: password hash failed")
		return errors.New("password reset failed")
	}

	userRepo := tx.GetUserRepository()
	user, err := userRepo.FindByID(ctx, userID) // Find user to update
	if err != nil {
		log.Error().Err(err).Msg("Failed to find user associated with token")
		s.handleRollback(tx, "Rollback ResetPassword: user find failed")
		return errors.New("password reset failed")
	}

	user.PasswordHash = hashedPassword                 // Update the hash
	if err := userRepo.Update(ctx, user); err != nil { // Save changes
		log.Error().Err(err).Msg("Failed to update user password")
		s.handleRollback(tx, "Rollback ResetPassword: user update failed")
		return errors.New("password reset failed")
	}

	// 7. Cabut Sesi (Log out all other sessions).
	// This deletes all their *refresh tokens*.
	personalTokenRepo := tx.GetPersonalTokenRepository()
	if err := personalTokenRepo.DeleteByUserID(ctx, userID); err != nil {
		log.Error().Err(err).Msg("Failed to revoke refresh tokens")
		s.handleRollback(tx, "Rollback ResetPassword: token revocation failed")
		return errors.New("password reset failed")
	}

	// 8. Cabut Token Reset (Delete the used token).
	if err := tokenRepo.DeleteByUserIDAndPurpose(ctx, userID, domain.TokenPurposePasswordReset); err != nil {
		log.Error().Err(err).Msg("Failed to delete used reset token")
		s.handleRollback(tx, "Rollback ResetPassword: reset token deletion failed")
		return errors.New("password reset failed")
	}

	// 9. Commit Transaksi.
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.New("password reset failed")
	}

	log.Info().Str("user_id", userID).Msg("Password reset successfully")
	return nil
}

// ResendVerificationEmail implements the logic from section 1.7.3.
func (s *authService) ResendVerificationEmail(ctx context.Context, userID string) error {
	log := s.logger.With().Str("method", "ResendVerificationEmail").Str("user_id", userID).Logger()

	// 1. Start Transaction
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return errors.New("request failed")
	}

	// Get transactional repositories
	userRepo := tx.GetUserRepository()
	tokenRepo := tx.GetUserTokenRepository()

	// 2. Find user and check if already verified
	user, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msg("User not found")
		} else {
			log.Error().Err(err).Msg("Failed to query user")
		}
		s.handleRollback(tx, "Rollback ResendVerification: user not found")
		return errors.New("request failed")
	}

	if user.IsVerified() { // Check using domain logic
		log.Warn().Msg("User email is already verified")
		s.handleRollback(tx, "Rollback ResendVerification: already verified")
		return errors.New("email already verified") //
	}

	// 3. Generate new token
	tokenString, err := generateSecureToken(32) // Use the same helper as ForgotPassword
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate secure token")
		s.handleRollback(tx, "Rollback ResendVerification: token generation failed")
		return errors.New("request failed")
	}

	// 4. Hash the token for storage
	tokenHash := hashToken(tokenString)
	tokenExpiry := time.Now().Add(15 * time.Minute) // 15-minute expiry

	// 5. Save (Upsert) the new token hash to the DB
	userToken := &domain.UserToken{
		UserID:    user.ID,
		Purpose:   domain.TokenPurposeVerification, // Use our new domain constant
		TokenHash: tokenHash,
		ExpiresAt: tokenExpiry,
	}
	if err := tokenRepo.Save(ctx, userToken); err != nil {
		log.Error().Err(err).Msg("Failed to upsert verification token")
		s.handleRollback(tx, "Rollback ResendVerification: token save failed")
		return errors.New("request failed")
	}

	// 6. Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.New("request failed")
	}

	// 7. Send email asynchronously
	go func() {
		emailCtx := context.Background()
		log.Info().Msg("Dispatching verification email")

		// The token format is `uid.tokenstring`
		fullToken := fmt.Sprintf("%s.%s", user.ID, tokenString)

		subject := "Verify Your Email Address"
		plainBody := fmt.Sprintf("Hi %s,\n\nPlease verify your email address using this token (it will expire in 15 minutes):\n\n%s\n\nIf you did not request this, please ignore this email.", user.Name, fullToken)
		htmlBody := fmt.Sprintf("<h1>Hi %s,</h1><p>Please verify your email address using this token (it will expire in 15 minutes):</p><h2>%s</h2><p>If you did not request this, please ignore this email.</p>", user.Name, fullToken)

		if err := s.mailService.Send(emailCtx, user.Email, subject, plainBody, htmlBody); err != nil {
			log.Error().Err(err).Msg("Failed to send verification email")
		}
	}()

	return nil
}

// VerifyEmail implements the logic from section 1.8.3.
func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	log := s.logger.With().Str("method", "VerifyEmail").Logger()

	// 1. Validate token format: "uid.tokenstring"
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return errors.New("invalid token format")
	}
	userID, tokenString := parts[0], parts[1]

	if userID == "" || tokenString == "" {
		return errors.New("invalid token format")
	}

	// 2. Start Transaction
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return errors.New("email verification failed")
	}

	// 3. Find token in DB
	tokenRepo := tx.GetUserTokenRepository()
	userToken, err := tokenRepo.FindByUserIDAndPurpose(ctx, userID, domain.TokenPurposeVerification)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msg("Email verification token not found in DB")
		} else {
			log.Error().Err(err).Msg("Failed to find token")
		}
		s.handleRollback(tx, "Rollback VerifyEmail: token not found")
		return errors.New("invalid or expired token") //
	}

	// 4. Hash the token from the user
	tokenHash := hashToken(tokenString)

	// 5. Validate the token hash and expiry
	if subtle.ConstantTimeCompare([]byte(userToken.TokenHash), []byte(tokenHash)) != 1 {
		log.Warn().Msg("Token hash mismatch")
		s.handleRollback(tx, "Rollback VerifyEmail: token hash mismatch")
		return errors.New("invalid or expired token") //
	}

	if time.Now().After(userToken.ExpiresAt) {
		log.Warn().Msg("Token expired")
		_ = tokenRepo.DeleteByUserIDAndPurpose(ctx, userID, domain.TokenPurposeVerification) // Clean up expired token
		s.handleRollback(tx, "Rollback VerifyEmail: token expired")
		return errors.New("invalid or expired token") //
	}

	// --- Token is valid ---

	// 6. Update user's verification status
	userRepo := tx.GetUserRepository()
	user, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find user associated with token")
		s.handleRollback(tx, "Rollback VerifyEmail: user find failed")
		return errors.New("email verification failed")
	}

	now := time.Now()
	user.EmailVerifiedAt = &now // Set the verification timestamp
	if err := userRepo.Update(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to update user verification status")
		s.handleRollback(tx, "Rollback VerifyEmail: user update failed")
		return errors.New("email verification failed")
	}

	// 7. Cabut Token (Delete the used token)
	if err := tokenRepo.DeleteByUserIDAndPurpose(ctx, userID, domain.TokenPurposeVerification); err != nil {
		log.Error().Err(err).Msg("Failed to delete used verification token")
		s.handleRollback(tx, "Rollback VerifyEmail: token deletion failed")
		return errors.New("email verification failed")
	}

	// 8. Commit Transaksi
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.New("email verification failed")
	}

	log.Info().Str("user_id", userID).Msg("Email verified successfully")
	return nil
}

// RequestEmailChange implements the logic from section 1.9.4.
func (s *authService) RequestEmailChange(ctx context.Context, userID string, newEmail string, currentPassword string) error {
	log := s.logger.With().Str("method", "RequestEmailChange").Str("user_id", userID).Logger()

	// 1. Start Transaction
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return errors.New("request failed")
	}

	// Get transactional repositories
	userRepo := tx.GetUserRepository()
	tokenRepo := tx.GetUserTokenRepository()

	// 2. Find user
	user, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msg("User not found")
		} else {
			log.Error().Err(err).Msg("Failed to query user")
		}
		s.handleRollback(tx, "Rollback RequestEmailChange: user not found")
		return errors.New("request failed")
	}

	// 3. Check if current email is verified
	if !user.IsVerified() {
		log.Warn().Msg("User's current email is not verified")
		s.handleRollback(tx, "Rollback RequestEmailChange: current email not verified")
		return errors.New("must verify current email before changing it")
	}

	// 4. Verifikasi Password
	if !s.hasher.Compare(user.PasswordHash, currentPassword) {
		log.Warn().Msg("Invalid current password provided")
		s.handleRollback(tx, "Rollback RequestEmailChange: invalid password")
		return errors.New("invalid password") //
	}

	// 5. Cek Duplikasi Email Baru
	_, err = userRepo.FindByEmail(ctx, newEmail)
	if err == nil {
		// An active user with this email already exists
		log.Warn().Str("new_email", newEmail).Msg("New email already in use")
		s.handleRollback(tx, "Rollback RequestEmailChange: new email in use")
		return errors.New("new email already in use") //
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// A database error occurred
		log.Error().Err(err).Msg("Failed to check new email")
		s.handleRollback(tx, "Rollback RequestEmailChange: db error")
		return errors.New("request failed")
	}

	// 6. Generate and save (Upsert) the new token
	tokenString, err := generateSecureToken(32)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate secure token")
		s.handleRollback(tx, "Rollback RequestEmailChange: token generation failed")
		return errors.New("request failed")
	}

	tokenHash := hashToken(tokenString)
	tokenExpiry := time.Now().Add(15 * time.Minute) // 15-minute expiry

	userToken := &domain.UserToken{
		UserID:    user.ID,
		Purpose:   domain.TokenPurposeEmailChange,
		TokenHash: tokenHash,
		ExpiresAt: tokenExpiry,
	}
	if err := tokenRepo.Save(ctx, userToken); err != nil {
		log.Error().Err(err).Msg("Failed to upsert email change token")
		s.handleRollback(tx, "Rollback RequestEmailChange: token save failed")
		return errors.New("request failed")
	}

	// 7. Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.New("request failed")
	}

	// 8. Send email asynchronously to the *new* email address
	go func() {
		emailCtx := context.Background()
		log.Info().Str("new_email", newEmail).Msg("Dispatching email change confirmation")

		// Create the token in the format: "uid.base64(new_email).tokenstring"
		// This matches the requirement to extract all 3 parts
		b64Email := base64.URLEncoding.EncodeToString([]byte(newEmail))
		fullToken := fmt.Sprintf("%s.%s.%s", user.ID, b64Email, tokenString)

		subject := "Confirm Your New Email Address"
		plainBody := fmt.Sprintf("Hi %s,\n\nPlease use this token to confirm your new email address (it will expire in 15 minutes):\n\n%s\n\nIf you did not request this, please ignore this email.", user.Name, fullToken)
		htmlBody := fmt.Sprintf("<h1>Hi %s,</h1><p>Please use this token to confirm your new email address (it will expire in 15 minutes):</p><h2>%s</h2><p>If you did not request this, please ignore this email.</p>", user.Name, fullToken)

		if err := s.mailService.Send(emailCtx, newEmail, subject, plainBody, htmlBody); err != nil {
			log.Error().Err(err).Str("new_email", newEmail).Msg("Failed to send email change confirmation")
		}
	}()

	return nil
}

// ConfirmEmailChange implements the logic from section 1.10.3.
func (s *authService) ConfirmEmailChange(ctx context.Context, token string) error {
	log := s.logger.With().Str("method", "ConfirmEmailChange").Logger()

	// 1. Validate token format: "uid.base64(new_email).tokenstring"
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("invalid token format")
	}
	userID, b64Email, tokenString := parts[0], parts[1], parts[2]

	// 2. Ekstrak Klaim
	emailBytes, err := base64.URLEncoding.DecodeString(b64Email)
	if err != nil || userID == "" || tokenString == "" {
		return errors.New("invalid token format")
	}
	newEmail := string(emailBytes)

	// 3. Hash Token
	tokenHash := hashToken(tokenString)

	// 4. Start Transaction
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return errors.New("email confirmation failed")
	}

	// 5. Find token in DB
	tokenRepo := tx.GetUserTokenRepository()
	userToken, err := tokenRepo.FindByUserIDAndPurpose(ctx, userID, domain.TokenPurposeEmailChange)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msg("Email change token not found in DB")
		} else {
			log.Error().Err(err).Msg("Failed to find token")
		}
		s.handleRollback(tx, "Rollback ConfirmEmailChange: token not found")
		return errors.New("invalid or expired token")
	}

	// 6. Validate Token (Hash and Expiry)
	if subtle.ConstantTimeCompare([]byte(userToken.TokenHash), []byte(tokenHash)) != 1 {
		log.Warn().Msg("Token hash mismatch")
		s.handleRollback(tx, "Rollback ConfirmEmailChange: token hash mismatch")
		return errors.New("invalid or expired token")
	}

	if time.Now().After(userToken.ExpiresAt) {
		log.Warn().Msg("Token expired")
		_ = tokenRepo.DeleteByUserIDAndPurpose(ctx, userID, domain.TokenPurposeEmailChange) // Clean up
		s.handleRollback(tx, "Rollback ConfirmEmailChange: token expired")
		return errors.New("invalid or expired token")
	}

	// --- Token is valid ---

	// 7. Get old user data (for notification)
	userRepo := tx.GetUserRepository()
	user, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find user associated with token")
		s.handleRollback(tx, "Rollback ConfirmEmailChange: user find failed")
		return errors.New("email confirmation failed")
	}
	oldEmail := user.Email

	// 8. Update Email
	now := time.Now()
	user.Email = newEmail
	user.EmailVerifiedAt = &now // Mark the new email as verified
	if err := userRepo.Update(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to update user email")
		s.handleRollback(tx, "Rollback ConfirmEmailChange: user update failed")
		return errors.New("email confirmation failed")
	}

	// 9. Cabut Token
	if err := tokenRepo.DeleteByUserIDAndPurpose(ctx, userID, domain.TokenPurposeEmailChange); err != nil {
		log.Error().Err(err).Msg("Failed to delete used email change token")
		s.handleRollback(tx, "Rollback ConfirmEmailChange: token deletion failed")
		return errors.New("email confirmation failed")
	}

	// 10. Commit Transaksi
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.New("email confirmation failed")
	}

	// 11. Send notification email to *old* address
	go func() {
		emailCtx := context.Background()
		log.Info().Str("old_email", oldEmail).Msg("Dispatching notification to old email")
		subject := "Your Email Address Has Been Changed"
		plainBody := fmt.Sprintf("Hi %s,\n\nThis is a notification that the email address for your account has been successfully changed to %s.\n\nIf you did not make this change, please contact support immediately.", user.Name, newEmail)
		htmlBody := fmt.Sprintf("<h1>Hi %s,</h1><p>This is a notification that the email address for your account has been successfully changed to <b>%s</b>.</p><p>If you did not make this change, please contact support immediately.</p>", user.Name, newEmail)

		if err := s.mailService.Send(emailCtx, oldEmail, subject, plainBody, htmlBody); err != nil {
			log.Error().Err(err).Str("old_email", oldEmail).Msg("Failed to send notification to old email")
		}
	}()

	return nil
}

// DisableAccount implements the logic from section 1.12.4.
func (s *authService) DisableAccount(ctx context.Context, userID string, currentPassword string) error {
	log := s.logger.With().Str("method", "DisableAccount").Str("user_id", userID).Logger()

	// 1. Start Transaction
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return errors.New("account disabling failed")
	}

	// Get transactional repositories
	userRepo := tx.GetUserRepository()
	personalTokenRepo := tx.GetPersonalTokenRepository()

	// 2. Find user (required for password check)
	user, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msg("User not found")
		} else {
			log.Error().Err(err).Msg("Failed to query user")
		}
		s.handleRollback(tx, "Rollback DisableAccount: user not found")
		return errors.New("account disabling failed")
	}

	// 3. Verifikasi Password
	if !s.hasher.Compare(user.PasswordHash, currentPassword) {
		log.Warn().Msg("Invalid current password provided for disable")
		s.handleRollback(tx, "Rollback DisableAccount: invalid password")
		return errors.New("invalid password") //
	}

	// --- Password is valid, proceed ---

	// 4. Nonaktifkan Akun (Set disabled_at)
	now := time.Now()
	user.DisabledAt = &now // Set the disabled timestamp
	if err := userRepo.Update(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to update user status to disabled")
		s.handleRollback(tx, "Rollback DisableAccount: user update failed")
		return errors.New("account disabling failed")
	}

	// 5. Cabut Sesi (Log out everywhere)
	// This deletes all their *refresh tokens*
	if err := personalTokenRepo.DeleteByUserID(ctx, userID); err != nil {
		log.Error().Err(err).Msg("Failed to revoke refresh tokens")
		s.handleRollback(tx, "Rollback DisableAccount: token revocation failed")
		return errors.New("account disabling failed")
	}

	// 6. Commit Transaksi
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.New("account disabling failed")
	}

	log.Info().Msg("User account disabled successfully")
	return nil
}

// RequestAccountReactivation sends a reactivation token
func (s *authService) RequestAccountReactivation(ctx context.Context, email string) error {
	log := s.logger.With().Str("method", "RequestAccountReactivation").Str("email", email).Logger()

	// Start a transaction
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return nil // Always return nil to prevent email enumeration
	}

	userRepo := tx.GetUserRepository()
	user, err := userRepo.FindByEmail(ctx, email)

	// We only proceed if the user exists AND is currently disabled
	if err != nil || !user.IsDisabled() {
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("Failed to query user by email")
		}
		// If user not found, or is found but *not* disabled, we do nothing.
		_ = tx.Commit() // Commit the empty transaction
		return nil      // Return ambiguous success
	}

	// --- User was found and is disabled, proceed with logic ---
	tokenRepo := tx.GetUserTokenRepository()

	// 1. Generate new token
	tokenString, err := generateSecureToken(32)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate secure token")
		s.handleRollback(tx, "Rollback Reactivation: token generation failed")
		return nil
	}

	// 2. Hash the token for storage
	tokenHash := hashToken(tokenString)
	tokenExpiry := time.Now().Add(15 * time.Minute) // 15-minute expiry

	// 3. Save (Upsert) the new token hash
	userToken := &domain.UserToken{
		UserID:    user.ID,
		Purpose:   domain.TokenPurposeAccountReactivation, // Use our new constant
		TokenHash: tokenHash,
		ExpiresAt: tokenExpiry,
	}
	if err := tokenRepo.Save(ctx, userToken); err != nil {
		log.Error().Err(err).Msg("Failed to upsert reactivation token")
		s.handleRollback(tx, "Rollback Reactivation: token save failed")
		return nil
	}

	// 4. Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return nil
	}

	// 5. Send email asynchronously
	go func() {
		emailCtx := context.Background()
		log.Info().Msg("Dispatching account reactivation email")
		// Use the same token format as password reset: uid.tokenstring
		fullToken := fmt.Sprintf("%s.%s", user.ID, tokenString)

		subject := "Account Reactivation Request"
		plainBody := fmt.Sprintf("Hi %s,\n\nWe received a request to reactivate your account. Use this token (it will expire in 15 minutes):\n\n%s\n\nIf you did not request this, please ignore this email.", user.Name, fullToken)
		htmlBody := fmt.Sprintf("<h1>Hi %s,</h1><p>We received a request to reactivate your account. Use this token (it will expire in 15 minutes):</p><h2>%s</h2><p>If you did not request this, please ignore this email.</p>", user.Name, fullToken)

		if err := s.mailService.Send(emailCtx, user.Email, subject, plainBody, htmlBody); err != nil {
			log.Error().Err(err).Msg("Failed to send reactivation email")
		}
	}()

	return nil
}

// ConfirmAccountReactivation validates a token and re-enables an account
func (s *authService) ConfirmAccountReactivation(ctx context.Context, token string) error {
	log := s.logger.With().Str("method", "ConfirmAccountReactivation").Logger()

	// 1. Validate token format: "uid.tokenstring"
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return errors.New("invalid token format")
	}
	userID, tokenString := parts[0], parts[1]
	if userID == "" || tokenString == "" {
		return errors.New("invalid token format")
	}

	// 2. Hash the provided token string
	tokenHash := hashToken(tokenString)

	// 3. Start Transaction
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return errors.New("account reactivation failed")
	}

	// 4. Find token in DB
	tokenRepo := tx.GetUserTokenRepository()
	userToken, err := tokenRepo.FindByUserIDAndPurpose(ctx, userID, domain.TokenPurposeAccountReactivation)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msg("Reactivation token not found in DB")
		} else {
			log.Error().Err(err).Msg("Failed to find token")
		}
		s.handleRollback(tx, "Rollback Reactivation: token not found")
		return errors.New("invalid or expired token")
	}

	// 5. Validate the token hash and expiry
	if subtle.ConstantTimeCompare([]byte(userToken.TokenHash), []byte(tokenHash)) != 1 {
		log.Warn().Msg("Token hash mismatch")
		s.handleRollback(tx, "Rollback Reactivation: token hash mismatch")
		return errors.New("invalid or expired token")
	}

	if time.Now().After(userToken.ExpiresAt) {
		log.Warn().Msg("Token expired")
		_ = tokenRepo.DeleteByUserIDAndPurpose(ctx, userID, domain.TokenPurposeAccountReactivation) // Clean up
		s.handleRollback(tx, "Rollback Reactivation: token expired")
		return errors.New("invalid or expired token")
	}

	// --- Token is valid ---

	// 6. Update user's status (set DisabledAt to NULL)
	userRepo := tx.GetUserRepository()
	user, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find user associated with token")
		s.handleRollback(tx, "Rollback Reactivation: user find failed")
		return errors.New("account reactivation failed")
	}

	user.DisabledAt = nil // This is the key step to re-enable the account
	if err := userRepo.Update(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to update user disabled status")
		s.handleRollback(tx, "Rollback Reactivation: user update failed")
		return errors.New("account reactivation failed")
	}

	// 7. Delete the used token
	if err := tokenRepo.DeleteByUserIDAndPurpose(ctx, userID, domain.TokenPurposeAccountReactivation); err != nil {
		log.Error().Err(err).Msg("Failed to delete used reactivation token")
		s.handleRollback(tx, "Rollback Reactivation: token deletion failed")
		return errors.New("account reactivation failed")
	}

	// 8. Commit Transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.New("account reactivation failed")
	}

	log.Info().Str("user_id", userID).Msg("Account reactivated successfully")
	return nil
}
