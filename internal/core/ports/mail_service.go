package ports

import "context"

// MailService defines the contract for sending emails.
// This is our "driven port".
type MailService interface {
	// Send sends an email.
	Send(
		ctx context.Context,
		to string,
		subject string,
		plainBody string,
		htmlBody string,
	) error
}
