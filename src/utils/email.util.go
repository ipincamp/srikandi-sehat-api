package utils

import (
	"fmt"
	"ipincamp/srikandi-sehat/config"
	"log"
)

func SendVerificationEmail(toEmail string, token string) error {
	baseURL := config.Get("APP_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:" + config.Get("API_PORT")
	}

	verificationLink := fmt.Sprintf("%s/api/auth/verify-email?token=%s", baseURL, token)

	// Simulasi pengiriman email dengan mencatatnya ke log
	log.Printf("===== SIMULASI PENGIRIMAN EMAIL =====")
	log.Printf("KE: %s", toEmail)
	log.Printf("SUBJEK: Verifikasi Akun Srikandi Sehat Anda")
	log.Printf("BODY: Silakan klik tautan berikut untuk memverifikasi email Anda:")
	log.Printf("%s", verificationLink)
	log.Printf("======================================")

	// Di produksi, jika pengiriman gagal, return error di sini
	return nil
}
