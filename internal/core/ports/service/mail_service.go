package service

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/ports/dto"
)

// MailService defines the contract for email delivery adapter handling various email notifications.
type MailService interface {
	// SendPasswordResetEmail sends a password reset OTP email to the user.
	SendPasswordResetEmail(ctx context.Context, dto dto.PasswordResetMailRequest) error

	// SendEmailVerificationEmail sends an email verification OTP to the user.
	SendEmailVerificationEmail(ctx context.Context, dto dto.EmailVerificationMailRequest) error

	// SendEmailChangeEmail sends an email change verification OTP to the new email address.
	SendEmailChangeEmail(ctx context.Context, dto dto.EmailChangeMailRequest) error

	// SendDeleteAccountEmail sends an account deletion notification to the user.
	SendDeleteAccountEmail(ctx context.Context, dto dto.AccountNotificationMailRequest) error

	// SendDisableAccountEmail sends an account disabled notification to the user.
	SendDisableAccountEmail(ctx context.Context, dto dto.AccountNotificationMailRequest) error
}
