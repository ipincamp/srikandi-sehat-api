package domain

import "time"

// User is the core domain model for a user.
// This struct represents the business entity, free of any
// database or transport layer details.
type User struct {

	// ID is the internal database identifier.
	ID string

	// Name of the user.
	Name string

	// Email is the user's login and contact email.
	Email string

	// EmailVerifiedAt is the timestamp when the user's email was verified.
	EmailVerifiedAt *time.Time

	// PasswordHash is the hashed password for authentication.
	PasswordHash string

	// DisabledAt is the timestamp when the user was disabled.
	DisabledAt *time.Time

	// CreatedAt timestamp.
	CreatedAt time.Time

	// UpdatedAt timestamp.
	UpdatedAt time.Time
}

// IsVerified returns true if the user's email is verified.
func (u *User) IsVerified() bool {
	return u.EmailVerifiedAt != nil
}

// IsDisabled returns true if the user account is disabled.
func (u *User) IsDisabled() bool {
	return u.DisabledAt != nil
}

// Logic that requires external dependencies (like password hashing)
// belongs in the 'service' layer.
