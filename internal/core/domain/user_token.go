package domain

import "time"

// TokenPurpose constants define the various purposes for user tokens.
const (
	TokenPurposeVerification    = "email_verification"
	TokenPurposePasswordReset   = "password_reset"
	TokenPurposeEmailChange     = "email_change"
	TokenPurposeAccountDeletion = "account_deletion"
)

// UserToken is the core domain model for a user token.
// This struct represents the business entity, free of any
// database or transport layer details.
type UserToken struct {

	// ID is the internal database identifier.
	UserID string

	// Purpose indicates the token's purpose (e.g., verification, password reset).
	Purpose string

	// TokenHash is the securely hashed token value.
	TokenHash string

	// ExpiresAt timestamp indicating when the token expires.
	ExpiresAt time.Time

	// CreatedAt timestamp.
	CreatedAt time.Time
}
