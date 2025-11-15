package service

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/ports/dto"
)

// MailService defines the contract for email delivery adapter handling various email notifications.
type MailService interface {
	// SendPasswordReset sends a password reset OTP email to the user.
	SendPasswordReset(ctx context.Context, dto dto.PasswordResetMailRequest) error

	// SendEmailVerification sends an email verification OTP to the user.
	SendEmailVerification(ctx context.Context, dto dto.EmailVerificationMailRequest) error

	// SendEmailChange sends an email change verification OTP to the new email address.
	SendEmailChange(ctx context.Context, dto dto.EmailChangeMailRequest) error

	// SendDeleteAccount sends an account deletion notification to the user.
	SendDeleteAccount(ctx context.Context, dto dto.AccountNotificationMailRequest) error

	// SendDisableAccount sends an account disabled notification to the user.
	SendDisableAccount(ctx context.Context, dto dto.AccountNotificationMailRequest) error
}
