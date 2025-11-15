package dto

// -- Lifecycle Request DTOs ---

// RegisterUserRequest represents the data required to register a new user.
type RegisterUserRequest struct {
	// Name of the user to be registered.
	Name string

	// Email of the user to be registered.
	Email string

	// Password for the new user account.
	Password string
}

// IdUserRequest represents a request that requires only a user ID.
type IdUserRequest struct {
	// UserID is the unique identifier of the user.
	ID string
}

// ForgotPasswordUserRequest represents the data required to initiate a forgot password flow.
type ForgotPasswordUserRequest struct {
	// Email of the user who forgot their password.
	Email string
}

// ResetPasswordUserRequest represents the data required to complete a password reset flow.
type ResetPasswordUserRequest struct {
	// Email of the user resetting their password.
	Email string

	// ResetToken is the token used to verify the password reset request.
	ResetToken string

	// NewPassword is the new password to be set for the user.
	NewPassword string
}

// -- Profile & Verification Request DTOs ---

// UpdateProfileUserRequest represents the data required to update a user's profile.
type UpdateProfileUserRequest struct {
	// UserID is the unique identifier of the user to be updated.
	IdUserRequest

	// NewName is the new name to update for the user.
	NewName *string
}

// ChangePasswordUserRequest represents the data required for a user to change their password.
type ChangePasswordUserRequest struct {
	// UserID is the unique identifier of the user changing their password.
	IdUserRequest

	// OldPassword is the current password of the user.
	OldPassword string

	// NewPassword is the new password to be set for the user.
	NewPassword string
}

// VerifyEmailUserRequest represents the data required to verify a user's email.
type VerifyEmailUserRequest struct {
	// UserID is the unique identifier of the user verifying their email.
	IdUserRequest

	// VerificationToken is the token used to verify the email.
	VerificationToken string
}
