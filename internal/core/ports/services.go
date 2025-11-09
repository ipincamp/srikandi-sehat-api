package ports

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
)

// AuthResponse is a Data Transfer Object (DTO) used to return
// tokens to the caller. It is not a core domain entity.
type AuthResponse struct {
	AccessToken  string
	RefreshToken string
}

// AuthService defines the "driving port" for authentication operations.
// This is the primary API for our authentication logic.
type AuthService interface {
	// Register creates a new user, hashes their password,
	// saves them, and returns a new set of auth tokens.
	Register(ctx context.Context, name, email, password string) (*AuthResponse, error)

	// Login validates user credentials and returns a new set of auth tokens.
	Login(ctx context.Context, email, password string) (*AuthResponse, error)

	// RefreshToken validates a refresh token and issues a new pair of tokens.
	// This allows for token rotation.
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)

	// Logout invalidates a refresh token (e.g., from a database).
	// This assumes refresh tokens are stateful (e.g., stored in a DB).
	Logout(ctx context.Context, refreshToken string) error
}

// UserService defines the "driving port" for user management operations.
// This handles business logic not directly related to authentication.
type UserService interface {
	// GetUserByID retrieves a user's public profile information.
	GetUserByID(ctx context.Context, uuid string) (*domain.User, error)

	// CreateUser creates a new user (e.g., for an admin panel).
	// This is distinct from 'Register' as it might have different logic
	// (e.g., doesn't auto-login, different validation) and returns
	// the user entity itself.
	CreateUser(ctx context.Context, name, email, password string) (*domain.User, error)

	// TODO: Add other methods like UpdateUser, DeleteUser, etc.
}
