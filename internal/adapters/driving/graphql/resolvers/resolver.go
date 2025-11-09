package resolvers

import (
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/rs/zerolog"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	authService ports.AuthService
	logger      zerolog.Logger
}

// It's the entry point for injecting core services into the adapter.
func NewResolver(authService ports.AuthService, logger zerolog.Logger) *Resolver {
	return &Resolver{
		authService: authService,
		logger:      logger,
	}
}
