package domain

import "time"

// PersonalToken is the core domain model for a personal access token.
// This struct represents the business entity, free of any
// database or transport layer details.
type PersonalToken struct {

	// ID is the internal database identifier.
	ID string

	// Token is the actual personal access token string.
	UserID string

	// Token is the actual personal access token string.
	ExpiresAt time.Time

	// CreatedAt is the timestamp when the token was created.
	CreatedAt time.Time
}
