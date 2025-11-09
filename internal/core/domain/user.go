package domain

import "time"

// User is the core domain model for a user.
// This struct represents the business entity, free of any
// database or transport layer details.
type User struct {
	// ID is the internal, auto-incrementing primary key.
	ID uint

	// UUID is the external-facing, unique identifier.
	UUID string

	// Name of the user.
	Name string

	// Email is the user's login and contact email.
	Email string

	// Password is the securely hashed password.
	// The domain model itself does not know *how* it's hashed,
	// only that it stores the resulting hash.
	Password string

	// CreatedAt timestamp.
	CreatedAt time.Time

	// UpdatedAt timestamp.
	UpdatedAt time.Time
}

// NOTE: We can add pure business logic methods here, for example:
//
// func (u *User) IsActive() bool {
//     // some logic
//     return true
// }
//
// Logic that requires external dependencies (like password hashing)
// belongs in the 'service' layer.
