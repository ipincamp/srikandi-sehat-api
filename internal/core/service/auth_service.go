package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/ipincamp/srikandi-sehat/pkg/config"
	"github.com/ipincamp/srikandi-sehat/pkg/password"
	"github.com/ipincamp/srikandi-sehat/pkg/token"

	// Import the new repository error
	db "github.com/ipincamp/srikandi-sehat/internal/adapters/driven/postgres"
)

// Compile-time check
var _ ports.AuthService = (*authService)(nil)

// authService implements the ports.AuthService interface.
// It depends on abstractions (interfaces) for its dependencies.
type authService struct {
	userRepo ports.UserRepository
	maker    token.Maker
	hasher   password.Hasher
	tokenCfg config.Token // Use the config struct for TTLs
}

// NewAuthService is the constructor for authService.
func NewAuthService(
	userRepo ports.UserRepository,
	maker token.Maker,
	hasher password.Hasher,
	tokenCfg config.Token,
) ports.AuthService {
	return &authService{
		userRepo: userRepo,
		maker:    maker,
		hasher:   hasher,
		tokenCfg: tokenCfg,
	}
}

// Register creates a new user, hashes their password,
// saves them, and returns a new set of auth tokens.
func (s *authService) Register(ctx context.Context, name, email, passwordStr string) (*ports.AuthResponse, error) {
	// 1. Check if user already exists
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		// User found, email is taken
		return nil, ErrEmailExists
	}
	if !errors.Is(err, db.ErrUserNotFound) {
		// A different, unexpected database error occurred
		return nil, err
	}

	// 2. Hash the password
	hashedPassword, err := s.hasher.Hash(passwordStr)
	if err != nil {
		return nil, err
	}

	// 3. Create the domain user
	user := &domain.User{
		UUID:     uuid.NewString(), // Generate a new public UUID
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		// ID, CreatedAt, UpdatedAt will be set by the repository
	}

	// 4. Save the user
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	// 5. Generate tokens
	return s.createTokenSet(user)
}

// Login validates user credentials and returns a new set of auth tokens.
func (s *authService) Login(ctx context.Context, email, passwordStr string) (*ports.AuthResponse, error) {
	// 1. Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// 2. Compare password
	if !s.hasher.Compare(user.Password, passwordStr) {
		return nil, ErrInvalidCredentials
	}

	// 3. Generate tokens
	return s.createTokenSet(user)
}

// RefreshToken validates a refresh token and issues a new pair of tokens.
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*ports.AuthResponse, error) {
	// 1. Validate the refresh token
	payload, err := s.maker.ValidateToken(refreshToken)
	if err != nil {
		if errors.Is(err, token.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	// 2. Check that it's actually a refresh token
	if payload.UseFor != token.UseForRefreshToken {
		return nil, ErrTokenUseMismatch
	}

	// 3. Find the user
	user, err := s.userRepo.FindByID(ctx, payload.UserID)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 4. Generate new tokens
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
