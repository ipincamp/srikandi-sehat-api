package utils

import (
	"fmt"
	"ipincamp/srikandi-sehat/config"
	"log"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

// SendEmail adalah fungsi inti pengirim email menggunakan GoMail.
func SendEmail(to, subject, htmlBody string) error {
	// 1. Ambil konfigurasi SMTP dari .env
	host := config.Get("SMTP_HOST")
	portStr := config.Get("SMTP_PORT")
	user := config.Get("SMTP_USER")
	pass := config.Get("SMTP_PASS")
	from := config.Get("SMTP_FROM")

	// 2. Konversi port string ke integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("ERROR: Invalid SMTP_PORT value: %v", err)
		return err
	}

	// 3. Buat pesan email
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	// 4. Buat Dialer (koneksi SMTP)
	// Catatan: Ini akan menggunakan TLS
	d := gomail.NewDialer(host, port, user, pass)

	// 5. Kirim email
	if err := d.DialAndSend(m); err != nil {
		ErrorLogger.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	InfoLogger.Printf("Email sent successfully to %s", to)
	return nil
}

// SendVerificationOTPEmail membuat template HTML dan memanggil SendEmail.
func SendVerificationOTPEmail(toEmail, otp string, expiresAt time.Time) error {
	subject := "Kode Verifikasi Akun Srikandi Sehat Anda"

	// --- 2. FORMAT WAKTU KEDALUWARSA ---
	loc, err := time.LoadLocation(config.Get("TIMEZONE"))
	if err != nil {
		// Jika gagal, fallback ke Waktu Server Lokal
		loc = time.Local
	}
	// Format waktu agar jelas: "14:35:02 WIB (30 Oktober 2025)"
	formattedTime := expiresAt.In(loc).Format("15:04:05 MST (2 January 2006)")

	// Buat template HTML sederhana untuk email
	htmlBody := fmt.Sprintf(`
	<div style="font-family: Arial, sans-serif; line-height: 1.6;">
		<h2>Verifikasi Akun Srikandi Sehat Anda</h2>
		<p>Terima kasih telah mendaftar. Silakan gunakan kode OTP berikut untuk menyelesaikan proses registrasi Anda:</p>
		<p style="font-size: 28px; font-weight: bold; letter-spacing: 4px; color: #333;">
			%s
		</p>
		<p style="color: #888;">
			Kode ini akan kedaluwarsa pada:<br>
			<strong style="color: #D9534F;">%s</strong>
		</p>
		<p>Jika Anda tidak merasa mendaftar, abaikan email ini.</p>
		<br>
		<p>Salam,</p>
		<p>Tim Srikandi Sehat</p>
	</div>
	`, otp, formattedTime)

	// Panggil pengirim email inti
	return SendEmail(toEmail, subject, htmlBody)
}
