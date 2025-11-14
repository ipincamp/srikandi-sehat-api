package dto

// PasswordResetMailRequest represents the data required to process a password reset request.
type PasswordResetMailRequest struct {
	// UserEmail is the email of the user requesting password reset.
	UserEmail string

	// ResetToken is the token used to reset the password.
	ResetToken string

	// Recipient Name
	Name string
}

// EmailVerificationMailRequest represents the data required to process an email verification request.
type EmailVerificationMailRequest struct {
	// UserEmail is the email of the user requesting email verification.
	UserEmail string

	// VerificationToken is the token used to verify the email.
	VerificationToken string

	// Recipient Name
	Name string
}

// EmailChangeMailRequest represents the data required to process an email change request.
type EmailChangeMailRequest struct {
	// OldEmail is the current email address of the user.
	OldEmail string

	// NewEmail is the new email address to be set for the user.
	NewEmail string

	// ChangeToken is the token used to verify the email change.
	ChangeToken string

	// Recipient Name
	Name string
}

// AccountNotificationMailRequest represents the data required to send user notifications.
type AccountNotificationMailRequest struct {
	// UserEmail is the email of the user to notify.
	UserEmail string

	// Recipient Name
	Name string
}
