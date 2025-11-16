package service

import (
	"context"
	"errors"
	"fmt"

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
	uow        ports.UnitOfWork
	hasher     password.Hasher
	tokenMaker token.Maker
	tokenCfg   config.Token
	logger     zerolog.Logger
}

// NewAuthService is the constructor for AuthService
func NewAuthService(
	uow ports.UnitOfWork,
	hasher password.Hasher,
	tokenMaker token.Maker,
	tokenCfg config.Token,
	logger zerolog.Logger,
) ports.AuthService {
	return &authService{
		uow:        uow,
		hasher:     hasher,
		tokenMaker: tokenMaker,
		tokenCfg:   tokenCfg,
		logger:     logger,
	}
}

// handleRollback is a helper to safely rollback a transaction and log failures
func (s *authService) handleRollback(tx ports.Transaction, logMsg string) {
	if err := tx.Rollback(); err != nil {
		s.logger.Error().Err(err).Msg(logMsg)
	}
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

	// Get transactional repositories
	userRepo := tx.GetUserRepository()

	// 2. Find the user by email (within tx)
	user, err := userRepo.FindByEmail(ctx, email)
	if err != nil {
		s.handleRollback(tx, "Rollback Login: find by email failed")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		s.logger.Error().Err(err).Msg("Failed to find user by email for login")
		return nil, errors.New("login failed")
	}

	// 3. Compare the password
	if !s.hasher.Compare(user.PasswordHash, password) {
		s.handleRollback(tx, "Rollback Login: invalid password")
		return nil, errors.New("invalid email or password")
	}

	// 4. Generate tokens and save JTI (within tx)
	authResponse, err := s.createTokenSet(ctx, tx, user)
	if err != nil {
		s.handleRollback(tx, "Rollback Login: failed to create token set")
		return nil, err
	}

	// 5. Commit Transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error().Err(err).Msg("Failed to commit transaction for Login")
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
