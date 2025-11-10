package mailgun

import (
	"context"
	"fmt"
	"net/smtp"
	"time"

	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/ipincamp/srikandi-sehat/pkg/config"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/rs/zerolog"
)

// Compile-time check
var _ ports.EmailService = (*mailService)(nil)

type mailService struct {
	cfg    config.Mail
	logger zerolog.Logger
	mg     mailgun.Mailgun // Mailgun client, bisa nil jika driver=smtp
}

// NewMailService membuat adapter email baru.
func NewMailService(cfg config.Mail, logger zerolog.Logger) ports.EmailService {
	svc := &mailService{
		cfg:    cfg,
		logger: logger,
	}

	if cfg.Driver == "mailgun" {
		svc.mg = mailgun.NewMailgun(cfg.MailgunDomain, cfg.MailgunAPIKey)
		logger.Info().Str("domain", cfg.MailgunDomain).Msg("Mailgun adapter initialized")
	} else {
		logger.Info().Str("host", cfg.Host).Int("port", cfg.Port).Msg("SMTP (Mailpit) adapter initialized")
	}

	return svc
}

// send adalah helper internal untuk mengirim email
func (s *mailService) send(ctx context.Context, to, subject, textBody string) error {
	from := fmt.Sprintf("%s <%s>", s.cfg.FromName, s.cfg.FromAddress)

	if s.cfg.Driver == "mailgun" {
		message := s.mg.NewMessage(from, subject, textBody, to)

		ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		resp, id, err := s.mg.Send(ctxTimeout, message)
		if err != nil {
			s.logger.Error().Err(err).Msg("Failed to send email via Mailgun")
			return err
		}
		s.logger.Info().Str("id", id).Str("resp", resp).Msg("Email sent via Mailgun")
		return nil
	}

	// --- Driver "smtp" (untuk Mailpit) ---
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	msg := []byte(
		"From: " + from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			textBody + "\r\n",
	)

	// Kirim email menggunakan net/smtp standar
	err := smtp.SendMail(addr, nil, s.cfg.FromAddress, []string{to}, msg)
	if err != nil {
		s.logger.Error().Err(err).Str("addr", addr).Msg("Failed to send email via SMTP (Mailpit)")
		return err
	}

	s.logger.Info().Str("to", to).Str("subject", subject).Msg("Email sent to SMTP (Mailpit)")
	return nil
}

// --- Implementasi Port ---

func (s *mailService) SendPasswordResetEmail(ctx context.Context, userEmail, name, otp string) error {
	subject := "Reset Password Akun Srikandi Sehat Anda"
	body := fmt.Sprintf(
		"Halo %s,\n\nAnda meminta untuk reset password.\n"+
			"Gunakan kode OTP ini untuk melanjutkan: %s\n\n"+
			"Kode ini akan kedaluwarsa dalam 15 menit.\n"+
			"Jika Anda tidak meminta ini, abaikan email ini.\n\nTerima kasih,\nTim Srikandi Sehat",
		name, otp,
	)
	return s.send(ctx, userEmail, subject, body)
}

func (s *mailService) SendEmailVerificationEmail(ctx context.Context, userEmail, name, otp string) error {
	subject := "Verifikasi Alamat Email Anda"
	body := fmt.Sprintf(
		"Halo %s,\n\nSelamat datang di Srikandi Sehat!\n"+
			"Untuk mengaktifkan akun Anda, masukkan kode OTP ini: %s\n\n"+
			"Kode ini akan kedaluwarsa dalam 1 jam.\n\nTerima kasih,\nTim Srikandi Sehat",
		name, otp,
	)
	return s.send(ctx, userEmail, subject, body)
}

func (s *mailService) SendEmailChangeEmail(ctx context.Context, oldEmail, newEmail, name, otp string) error {
	subject := "Konfirmasi Perubahan Email Srikandi Sehat"
	body := fmt.Sprintf(
		"Halo %s,\n\nKami menerima permintaan untuk mengubah email Anda dari %s ke %s.\n"+
			"Masukkan kode OTP ini untuk mengonfirmasi perubahan: %s\n\n"+
			"Jika Anda tidak meminta ini, abaikan email ini.\n\nTerima kasih,\nTim Srikandi Sehat",
		name, oldEmail, newEmail, otp,
	)
	// TODO: Idealnya, kirim ke email LAMA, bukan yang baru.
	// Tapi untuk blueprint ini, kita kirim ke email lama.
	return s.send(ctx, oldEmail, subject, body)
}

func (s *mailService) SendDeleteAccountEmail(ctx context.Context, userEmail, name string) error {
	subject := "Akun Srikandi Sehat Anda Telah Dihapus"
	body := fmt.Sprintf(
		"Halo %s,\n\nAkun Anda yang terdaftar dengan email %s telah berhasil dihapus.\n\n"+
			"Kami sedih melihat Anda pergi. Jika ini adalah kesalahan, silakan hubungi dukungan.\n\n"+
			"Terima kasih,\nTim Srikandi Sehat",
		name, userEmail,
	)
	return s.send(ctx, userEmail, subject, body)
}
