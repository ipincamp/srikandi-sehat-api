package service

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/ports/dto"
)

// AuthService defines the driving port for authentication operations including identity verification and token management.
type AuthService interface {
	// Login validates user credentials and returns authentication tokens on success.
	Login(ctx context.Context, dto dto.LoginAuthRequest) (*dto.TokenAuthResponse, error)

	// RefreshToken validates a refresh token and issues a new token pair.
	RefreshToken(ctx context.Context, dto dto.RefreshTokenAuthRequest) (*dto.TokenAuthResponse, error)

	// Logout invalidates a refresh token.
	Logout(ctx context.Context, dto dto.LogoutAuthRequest) error
}
