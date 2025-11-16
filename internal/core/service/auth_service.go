package service

import (
	"context"
	"errors"

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
	userRepo          ports.UserRepository
	personalTokenRepo ports.PersonalTokenRepository
	hasher            password.Hasher
	tokenMaker        token.Maker
	tokenCfg          config.Token
	logger            zerolog.Logger
}

// NewAuthService is the constructor for AuthService
func NewAuthService(
	userRepo ports.UserRepository,
	personalTokenRepo ports.PersonalTokenRepository,
	hasher password.Hasher,
	tokenMaker token.Maker,
	tokenCfg config.Token,
	logger zerolog.Logger,
) ports.AuthService {
	return &authService{
		userRepo:          userRepo,
		personalTokenRepo: personalTokenRepo,
		hasher:            hasher,
		tokenMaker:        tokenMaker,
		tokenCfg:          tokenCfg,
		logger:            logger,
	}
}

// Register implements the registration logic
func (s *authService) Register(ctx context.Context, name, email, password string) (*domain.AuthResponse, error) {
	// 1. Check if user already exists
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		// User found, email is already taken
		return nil, errors.New("email already in use")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// A different database error occurred
		s.logger.Error().Err(err).Msg("Failed to check user by email")
		return nil, errors.New("registration failed")
	}
	// If err is gorm.ErrRecordNotFound, we can proceed

	// 2. Hash the password
	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to hash password")
		return nil, errors.New("registration failed")
	}

	// 3. Create the domain user
	newUser := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
		// ID, CreatedAt, etc., will be set by the database (see migration)
	}

	// 4. Save the user
	// Note: Your docs mention a Unit of Work (UoW). For simplicity, I am calling
	// the repo directly. In a real scenario, you'd start a UoW transaction here.
	if err := s.userRepo.Save(ctx, newUser); err != nil {
		s.logger.Error().Err(err).Msg("Failed to save new user")
		return nil, errors.New("registration failed")
	}

	// 5. Fetch the newly created user to get the ID
	createdUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to fetch newly created user")
		return nil, errors.New("registration failed")
	}

	// 6. Generate tokens
	return s.createTokenSet(ctx, createdUser)
}

// Login implements the login logic
func (s *authService) Login(ctx context.Context, email, password string) (*domain.AuthResponse, error) {
	// 1. Find the user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		s.logger.Error().Err(err).Msg("Failed to find user by email for login")
		return nil, errors.New("login failed")
	}

	// 2. Compare the password
	if !s.hasher.Compare(user.PasswordHash, password) {
		return nil, errors.New("invalid email or password")
	}

	// 3. Generate tokens
	return s.createTokenSet(ctx, user)
}

// createTokenSet is a helper to generate both access and refresh tokens
func (s *authService) createTokenSet(ctx context.Context, user *domain.User) (*domain.AuthResponse, error) {
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

	// Save the refresh token JTI to the database
	personalToken := &domain.PersonalToken{
		ID:        refreshPayload.JTI,
		UserID:    user.ID,
		ExpiresAt: refreshPayload.ExpiresAt,
	}

	if err := s.personalTokenRepo.Save(ctx, personalToken); err != nil {
		s.logger.Error().Err(err).Msg("Failed to save personal token JTI")
		return nil, errors.New("token generation failed")
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken validates an old refresh token and issues a new pair.
func (s *authService) RefreshToken(ctx context.Context, tokenString string) (*domain.AuthResponse, error) {
	// 1. Validate the token string
	payload, err := s.tokenMaker.ValidateToken(tokenString)
	if err != nil {
		return nil, errors.New("invalid refresh token") // Catches expired, malformed, etc.
	}

	// 2. Check its purpose
	if payload.UseFor != token.UseForRefreshToken {
		return nil, errors.New("invalid token purpose")
	}

	// 3. Check if the JTI is in our database
	if _, err := s.personalTokenRepo.FindByID(ctx, payload.JTI); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn().Str("jti", payload.JTI).Msg("Refresh token JTI not found in DB (already used or invalid)")
			return nil, errors.New("invalid refresh token")
		}
		s.logger.Error().Err(err).Msg("Failed to find personal token by ID")
		return nil, errors.New("token refresh failed")
	}

	// 4. JTI is valid. Revoke it (Token Rotation).
	if err := s.personalTokenRepo.Delete(ctx, payload.JTI); err != nil {
		s.logger.Error().Err(err).Msg("Failed to delete old JTI during refresh")
		return nil, errors.New("token refresh failed")
	}

	// 5. Fetch the user
	user, err := s.userRepo.FindByID(ctx, payload.UserID)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", payload.UserID).Msg("User for refresh token not found")
		return nil, errors.New("invalid refresh token")
	}

	// 6. Create a new token set (which saves the new JTI)
	return s.createTokenSet(ctx, user)
}

// Logout revokes a specific refresh token.
func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	// 1. Validate the token
	payload, err := s.tokenMaker.ValidateToken(refreshToken)
	if err != nil {
		// If token is already invalid, we can just return success
		s.logger.Warn().Err(err).Msg("Logout attempt with invalid token")
		return nil
	}

	// 2. Check its purpose
	if payload.UseFor != token.UseForRefreshToken {
		s.logger.Warn().Msg("Logout attempt with non-refresh token")
		return nil // Still success, client is just confused
	}

	// 3. Delete the JTI from the database
	if err := s.personalTokenRepo.Delete(ctx, payload.JTI); err != nil {
		// Log the error, but don't fail. The client wants to be logged out.
		s.logger.Error().Err(err).Str("jti", payload.JTI).Msg("Failed to delete JTI on logout")
	}

	return nil
}
