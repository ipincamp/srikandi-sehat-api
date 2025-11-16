package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrTokenExpired is a custom error returned when validation fails
// specifically because the token is past its expiration time.
// Defining this as a package-level variable allows upstream services
// to check for it using `errors.Is(err, token.ErrTokenExpired)` and
// provide a specific user response (e.g., HTTP 401) instead of a
// generic HTTP 500.
var ErrTokenExpired = errors.New("token has expired")

// UseFor constants define the valid 'use-for' claims.
// Using constants prevents "magic strings" in the codebase,
// reducing typos and making intent clear.
const (
	UseForAccessToken  = "access_token"
	UseForRefreshToken = "refresh_token"
)

// Payload contains data (claims) stored in the token.
// This struct is the "public" data structure for the token package,
// used by both the interface and the implementations.
type Payload struct {
	JTI       string    `json:"jti"`
	UserID    string    `json:"uid"`
	RoleID    string    `json:"rid"`
	UseFor    string    `json:"for"` // e.g., "access_token" or "refresh_token"
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
}

// NewPayload creates a new token payload.
// This factory function ensures all necessary fields are validated
// before a payload is created.
func NewPayload(userID, roleID, useFor string, duration time.Duration) (*Payload, error) {
	if userID == "" {
		return nil, errors.New("userID cannot be empty")
	}
	if useFor == "" {
		return nil, errors.New("useFor cannot be empty")
	}

	now := time.Now().UTC() // Use UTC for consistency
	return &Payload{
		JTI:       uuid.NewString(),
		UserID:    userID,
		RoleID:    roleID, // Can be empty (e.g., for refresh_token)
		UseFor:    useFor,
		IssuedAt:  now,
		ExpiresAt: now.Add(duration),
	}, nil
}

// Validate checks if the token payload is still valid.
// PASETO/JWT libraries often do this automatically, but it's
// good practice to have a method on the payload itself.
func (p *Payload) Validate() error {
	if time.Now().UTC().After(p.ExpiresAt) {
		return ErrTokenExpired
	}
	return nil
}
