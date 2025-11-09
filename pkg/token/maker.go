package token

import (
	"time"
)

// Maker defines contracts for all token operations.
// Like the Hasher, this is a crucial abstraction (DIP).
// Our application services will depend on `token.Maker`, not on
// `token.PasetoMaker` or `token.JwtMaker`.
// This allows us to swap out the entire token technology (e.g., from
// PASETO to JWT) without modifying any business logic, just the
// wiring in main.go.
type Maker interface {
	// CreateToken creates a new token with a custom payload.
	// It returns the token string and the payload struct used to create it.
	CreateToken(userID, roleID, useFor string, duration time.Duration) (string, *Payload, error)

	// ValidateToken verifies the token and returns its payload.
	// This is the inverse of CreateToken.
	ValidateToken(token string) (*Payload, error)
}
