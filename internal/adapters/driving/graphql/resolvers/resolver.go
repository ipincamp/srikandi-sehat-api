package resolvers

import (
	"errors"

	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/rs/zerolog"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// This context key is used by the auth middleware to pass the
// authenticated user's UUID to the resolver chain.
type contextKey string

const AuthUserUUIDKey contextKey = "authUserUUID"

var ErrNotAuthenticated = errors.New("not authenticated")

type Resolver struct {
	authService ports.AuthService
	userService ports.UserService
	logger      zerolog.Logger
}

// It's the entry point for injecting core services into the adapter.
func NewResolver(
	authService ports.AuthService,
	userService ports.UserService,
	logger zerolog.Logger,
) *Resolver {
	return &Resolver{
		authService: authService,
		userService: userService,
		logger:      logger,
	}
}
