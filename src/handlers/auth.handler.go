package handlers

import (
	"errors"
	"fmt"
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"

	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Helper untuk validasi domain
func isDomainAllowed(email string) bool {
	domainsStr := config.Get("ALLOWED_EMAIL_DOMAINS")
	if domainsStr == "" {
		// Jika tidak diset, izinkan semua (default aman)
		return true
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false // Format email tidak valid
	}
	domain := parts[1]

	allowedDomainsList := strings.Split(domainsStr, ",")
	allowedDomainsMap := make(map[string]bool)
	for _, d := range allowedDomainsList {
		allowedDomainsMap[strings.TrimSpace(d)] = true
	}

	return allowedDomainsMap[domain]
}

func Register(c *fiber.Ctx) error {
	input := c.Locals("request_body").(*dto.RegisterRequest)

	// Validasi Domain Email
	if !isDomainAllowed(input.Email) {
		return utils.SendError(c, fiber.StatusUnprocessableEntity, "Registrasi dari domain email ini tidak diizinkan.")
	}

	var existingUser models.User
	err := database.DB.First(&existingUser, "email = ?", input.Email).Error

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.SendError(c, fiber.StatusConflict, "Email sudah terdaftar")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.ErrorLogger.Printf("Gagal hash password untuk %s: %v", input.Email, err)
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memproses akun")
	}

	defaultRole, err := utils.GetRoleByName(string(constants.UserRole))
	if err != nil {
		utils.ErrorLogger.Printf("FATAL: Default role '%s' not found", constants.UserRole)
		return utils.SendError(c, fiber.StatusInternalServerError, "Konfigurasi server error")
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	// 4. Lakukan transaksi Database
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		// Gunakan defaultRole yang sudah diambil
		if err := tx.Model(&user).Association("Roles").Append(&defaultRole); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		utils.ErrorLogger.Printf("Gagal membuat user %s di db: %v", input.Email, err)
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal membuat akun")
	}

	user.Roles = []*models.Role{&defaultRole}
	user.Profile = models.Profile{} // Profile baru, ID-nya 0

	token, err := utils.GenerateJWT(user)
	if err != nil {
		utils.ErrorLogger.Printf("Gagal membuat JWT untuk user baru %s: %v", user.Email, err)
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memproses sesi login")
	}

	utils.AddEmailToRegistrationFilter(user.Email)

	responseData := dto.UserResponseJson(user, token)

	return utils.SendSuccess(c, fiber.StatusCreated, "Registrasi sukses! Silakan verifikasi email Anda untuk mendapatkan akses penuh.", responseData)
}

func VerifyOTP(c *fiber.Ctx) error {
	userUUID, ok := c.Locals("user_id").(string)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}
	input := c.Locals("request_body").(*dto.VerifyOTPRequest)

	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	// Cek 1: Apakah sudah terverifikasi?
	if user.EmailVerifiedAt.Valid {
		return utils.SendError(c, fiber.StatusConflict, "Email Anda sudah terverifikasi.")
	}

	// Cek 2: Apakah token/waktu kedaluwarsa valid?
	if !user.VerificationToken.Valid || !user.VerificationExpiresAt.Valid {
		return utils.SendError(c, fiber.StatusForbidden, "Tidak ada OTP yang aktif. Silakan minta kirim ulang.")
	}

	// Cek 3: Apakah kedaluwarsa?
	if time.Now().After(user.VerificationExpiresAt.Time) {
		return utils.SendError(c, fiber.StatusForbidden, "Kode OTP telah kedaluwarsa. Silakan minta kirim ulang.")
	}

	// Cek 4: Apakah OTP cocok?
	if user.VerificationToken.String != input.OTP {
		return utils.SendError(c, fiber.StatusUnauthorized, "Kode OTP salah.")
	}

	// Sukses! Verifikasi pengguna.
	updates := map[string]interface{}{
		"email_verified_at":       time.Now(),
		"verification_token":      nil,
		"verification_expires_at": nil,
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memverifikasi akun")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Email berhasil diverifikasi! Anda sekarang memiliki akses penuh.", nil)
}

func ResendVerification(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)

	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	if user.EmailVerifiedAt.Valid {
		return utils.SendError(c, fiber.StatusConflict, "Email sudah terverifikasi.")
	}

	if user.LastOTPSentAt.Valid {
		// Tentukan kapan user boleh request lagi
		nextAllowedTime := user.LastOTPSentAt.Time.Add(15 * time.Minute)

		if time.Now().Before(nextAllowedTime) {
			// Jika sekarang SEBELUM waktu yang diizinkan, tolak request
			remaining := time.Until(nextAllowedTime)

			// Format sisa waktu agar lebih ramah
			remainingMinutes := int(remaining.Minutes())
			remainingSeconds := int(remaining.Seconds()) % 60

			errMsg := fmt.Sprintf(
				"Harap tunggu %d menit %d detik lagi sebelum meminta kode baru.",
				remainingMinutes,
				remainingSeconds,
			)
			// Kirim status 429 Too Many Requests
			return utils.SendError(c, fiber.StatusTooManyRequests, errMsg)
		}
	}

	// Buat 6-digit OTP, kedaluwarsa 10 Menit
	verificationToken, err := utils.GenerateOTP(6)
	if err != nil {
		utils.ErrorLogger.Printf("Gagal membuat OTP (resend): %v", err)
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memproses permintaan")
	}
	verificationExpires := time.Now().Add(10 * time.Minute) // 10 Menit

	updates := map[string]interface{}{
		"verification_token":      verificationToken,
		"verification_expires_at": verificationExpires,
		"last_otp_sent_at":        time.Now(),
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memperbarui token")
	}

	// Kirim email (lagi)
	if err := utils.SendVerificationOTPEmail(user.Email, verificationToken, verificationExpires); err != nil {
		utils.ErrorLogger.Printf("Gagal mengirim ulang OTP ke %s: %v", user.Email, err)
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal mengirim email verifikasi.")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Kode OTP baru telah dikirim ke email Anda.", nil)
}

func Login(c *fiber.Ctx) error {
	input := c.Locals("request_body").(*dto.LoginRequest)

	var user models.User
	err := database.DB.Preload("Roles").Preload("Profile").First(&user, "email = ?", input.Email).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.SendError(c, fiber.StatusUnauthorized, "Invalid credentials")
	}
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Database query error")
	}

	match, err := utils.CheckPasswordHash(input.Password, user.Password)
	if err != nil || !match {
		return utils.SendError(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to generate JWT token")
	}

	go utils.AddUserToFrequentLoginFilter(user)
	responseData := dto.UserResponseJson(user, token)

	return utils.SendSuccess(c, fiber.StatusOK, "Login successful", responseData)
}

func Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenString == "" {
		return utils.SendError(c, fiber.StatusUnauthorized, "JWT token not found or invalid format")
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Malformed token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid expiration time in token")
	}
	expiresAt := time.Unix(int64(exp), 0)

	if time.Now().After(expiresAt) {
		return utils.SendSuccess(c, fiber.StatusOK, "Token already expired", nil)
	}

	// invalidToken := models.InvalidToken{
	// 	Token:     tokenString,
	// 	ExpiresAt: expiresAt,
	// }

	// if err := database.DB.Create(&invalidToken).Error; err != nil {
	// 	if strings.Contains(err.Error(), "Duplicate entry") {
	// 		// Token is already blocklisted
	// 	} else {
	// 		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to invalidate token")
	// 	}
	// }

	// remainingDuration := time.Until(expiresAt)
	// utils.AddToBlocklistCache(tokenString, remainingDuration)

	return utils.SendSuccess(c, fiber.StatusOK, "Logout successful", nil)
}
