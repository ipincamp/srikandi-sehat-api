package utils

import (
	"context"
	"errors"
	"fmt"
	"ipincamp/srikandi-sehat/config"

	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

func VerifyGoogleIDToken(idToken string) (*oauth2.Tokeninfo, error) {
	ctx := context.Background()
	oauth2Service, err := oauth2.NewService(ctx, option.WithoutAuthentication())
	if err != nil {
		ErrorLogger.Printf("Gagal membuat service oauth2: %v", err)
		return nil, errors.New("gagal memvalidasi token")
	}

	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)

	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		ErrorLogger.Printf("Gagal memverifikasi token: %v", err)
		return nil, errors.New("token tidak valid atau kedaluwarsa")
	}

	// Verifikasi Audience (Client ID)
	googleClientID := config.Get("GOOGLE_CLIENT_ID")
	if tokenInfo.Audience != googleClientID {
		errMsg := fmt.Sprintf("Token audience tidak cocok. Diharapkan: %s, Didapat: %s", googleClientID, tokenInfo.Audience)
		ErrorLogger.Println(errMsg)
		return nil, errors.New("token tidak valid (audience mismatch)")
	}

	return tokenInfo, nil
}
