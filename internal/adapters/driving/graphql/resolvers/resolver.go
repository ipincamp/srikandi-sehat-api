package resolvers

import "github.com/ipincamp/srikandi-sehat/internal/core/ports"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	authService ports.AuthService
}

// It's the entry point for injecting core services into the adapter.
func NewResolver(authService ports.AuthService) *Resolver {
	// When services are added, they would be injected here.
	// Example:
	// return &Resolver{userService: us}
	return &Resolver{
		authService: authService,
	}
}
