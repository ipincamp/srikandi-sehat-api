package utils

import (
	"context"
	"errors"
	"ipincamp/srikandi-sehat/config"
	"strings"

	"google.golang.org/api/idtoken"
)

type GoogleTokenPayload struct {
	Email string
	Name  string
	Sub   string // Ini adalah Google User ID
}

// VerifyGoogleIDToken memverifikasi idToken yang diberikan dan mengembalikan info token
func VerifyGoogleIDToken(idToken string) (*GoogleTokenPayload, error) {
	ctx := context.Background()
	googleClientID := config.Get("GOOGLE_CLIENT_ID")

	// Validasi token menggunakan library idtoken
	payload, err := idtoken.Validate(ctx, idToken, googleClientID)
	if err != nil {
		ErrorLogger.Printf("Gagal memverifikasi id token google: %v", err)
		return nil, errors.New("token tidak valid atau kedaluwarsa")
	}

	// Payload valid, ekstrak data yang kita butuhkan dari Claims
	claims := payload.Claims

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return nil, errors.New("token tidak mengandung email")
	}

	name, ok := claims["name"].(string)
	if !ok || name == "" {
		// Kadang 'name' tidak ada, fallback ke bagian email
		name = strings.Split(email, "@")[0]
	}

	return &GoogleTokenPayload{
		Email: email,
		Name:  name,
		Sub:   payload.Subject, // payload.Subject adalah Google User ID (Sub = Subject)
	}, nil
}
