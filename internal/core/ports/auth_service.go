package ports

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
)

// AuthService defines the business logic for authentication.
type AuthService interface {
	Register(ctx context.Context, name, email, password string) (*domain.AuthResponse, error)
	Login(ctx context.Context, email, password string) (*domain.AuthResponse, error)
	RefreshToken(ctx context.Context, tokenString string) (*domain.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}
