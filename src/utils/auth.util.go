package utils

import (
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/models"
	"log"
	"strconv"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
)

var argon2Params = &argon2id.Params{
	Memory:      64 * 1024, // 64 MB
	Iterations:  3,         // Recommended minimum is 3
	Parallelism: 2,         // Use 2 threads for efficiency
	SaltLength:  16,        // 16 bytes salt
	KeyLength:   32,        // 32 bytes key
}

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2Params)
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func GenerateJWT(user models.User) (string, error) {
	secretKey := config.Get("JWT_SECRET")

	expHoursStr := config.Get("JWT_EXPIRATION_HOURS")
	expHours, err := strconv.Atoi(expHoursStr)
	if err != nil {
		expHours = 24
	}

	claims := jwt.MapClaims{
		"usr": user.UUID,
		"exp": time.Now().Add(time.Hour * time.Duration(expHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func CleanupExpiredTokens() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Running expired token cleanup...")
		now := time.Now()
		result := database.DB.Where("expires_at < ?", now).Delete(&models.InvalidToken{})
		if result.Error != nil {
			log.Printf("Failed to clean up expired tokens: %v", result.Error)
		} else {
			log.Printf("%d expired tokens have been deleted.", result.RowsAffected)
		}
	}
}
