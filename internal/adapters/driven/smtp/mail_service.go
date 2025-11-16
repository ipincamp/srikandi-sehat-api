package smtp

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/ipincamp/srikandi-sehat/pkg/config"
	"github.com/rs/zerolog"
)

// mailService is the concrete implementation of the MailService port.
type mailService struct {
	cfg    config.Mail
	logger zerolog.Logger
}

// NewSMTPService creates a new mail service adapter.
func NewSMTPService(cfg config.Mail, logger zerolog.Logger) ports.MailService {
	return &mailService{
		cfg:    cfg,
		logger: logger,
	}
}

// Send implements the MailService interface.
func (s *mailService) Send(
	ctx context.Context,
	to string,
	subject string,
	plainBody string,
	htmlBody string,
) error {
	// 1. Check context cancellation before proceeding.
	if err := ctx.Err(); err != nil {
		s.logger.Warn().Err(err).Msg("Context cancelled before email send")
		return err
	}

	// 2. Conditionally set up the authentication mechanism.
	var auth smtp.Auth

	// We only set up authentication if a username is configured.
	// This allows Mailpit (no user) to work, while Google SMTP (has user)
	// will still authenticate properly.
	if s.cfg.SMTPUser != "" {
		s.logger.Debug().Msg("SMTPUser is set, using PlainAuth.")
		auth = smtp.PlainAuth(
			"",                 // identity (usually empty)
			s.cfg.SMTPUser,     // username
			s.cfg.SMTPPassword, // password
			s.cfg.Host,         // host
		)
	} else {
		s.logger.Debug().Msg("SMTPUser is empty, using nil Auth (Mailpit).")
	}

	// 3. Construct the email message with MIME headers.
	msg, err := s.buildMIMEMessage(to, subject, plainBody, htmlBody)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to build MIME message")
		return err
	}

	// 4. Set the server address (e.g., "localhost:1025").
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	// 5. Send the email.
	// When 'auth' is nil, SendMail does not attempt to authenticate.
	err = smtp.SendMail(
		addr,              // SMTP server address
		auth,              // Authentication (nil for Mailpit)
		s.cfg.FromAddress, // FROM address
		[]string{to},      // TO address(es)
		[]byte(msg),       // The full message body
	)

	if err != nil {
		s.logger.Error().Err(err).Str("to", to).Msg("Failed to send email via SMTP")
		return err
	}

	s.logger.Info().Str("to", to).Str("subject", subject).Msg("Email sent successfully")
	return nil
}

// buildMIMEMessage is a helper to construct a multipart email.
func (s *mailService) buildMIMEMessage(to, subject, plainBody, htmlBody string) (string, error) {
	var msg strings.Builder

	// Set headers
	fromHeader := fmt.Sprintf("%s <%s>", s.cfg.FromName, s.cfg.FromAddress)
	msg.WriteString(fmt.Sprintf("From: %s\r\n", fromHeader))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")

	// --- Create a multipart/alternative message ---
	// This allows email clients to choose between plain text and HTML.
	boundary := "boundary-mixed-part" // A unique boundary string
	msg.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n\r\n", boundary))

	// -- Plain Text Part --
	msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msg.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	msg.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	msg.WriteString(plainBody)
	msg.WriteString("\r\n")

	// -- HTML Part --
	if htmlBody != "" {
		msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		msg.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
		msg.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
		msg.WriteString(htmlBody)
		msg.WriteString("\r\n")
	}

	// -- End boundary --
	msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return msg.String(), nil
}
