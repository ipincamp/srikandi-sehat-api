package handlers

import (
	"database/sql"
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
		utils.AuthLogger.Printf("Registration failed (domain not allowed): %s", input.Email)
		return utils.SendError(c, fiber.StatusUnprocessableEntity, "Registrasi dari domain email ini tidak diizinkan.")
	}

	// Cek apakah email sudah ada di tabel users
	var existingUser models.User
	err := database.DB.First(&existingUser, "email = ?", input.Email).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.AuthLogger.Printf("Registration failed (email exists): %s", input.Email)
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

	// Buat user baru dan provider auth 'local' dalam satu transaksi
	var user models.User
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Buat User
		user = models.User{
			Name:  input.Name,
			Email: input.Email,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// 2. Buat Auth Provider
		authProvider := models.UserAuthProvider{
			UserID:   user.ID,
			Provider: "local",
			Password: sql.NullString{String: hashedPassword, Valid: true},
		}
		if err := tx.Create(&authProvider).Error; err != nil {
			return err
		}

		// 3. Tetapkan Role
		if err := tx.Model(&user).Association("Roles").Append(&defaultRole); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		utils.ErrorLogger.Printf("Gagal membuat user %s di db: %v", input.Email, err)
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal membuat akun")
	}

	// Siapkan response
	user.Roles = []*models.Role{&defaultRole}
	user.Profile = models.Profile{}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		utils.ErrorLogger.Printf("Gagal membuat JWT untuk user baru %s: %v", user.Email, err)
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memproses sesi login")
	}

	utils.AddEmailToRegistrationFilter(user.Email)

	responseData := dto.UserResponseJson(user, token)

	utils.AuthLogger.Printf("User registered successfully: %s (UUID: %s)", responseData.Email, responseData.ID)
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
		utils.AuthLogger.Printf("OTP verification skipped (already verified): %s", userUUID)
		return utils.SendError(c, fiber.StatusConflict, "Email Anda sudah terverifikasi.")
	}

	// Cek 2: Apakah token/waktu kedaluwarsa valid?
	if !user.VerificationToken.Valid || !user.VerificationExpiresAt.Valid {
		utils.AuthLogger.Printf("OTP verification failed (no active OTP): %s", userUUID)
		return utils.SendError(c, fiber.StatusForbidden, "Tidak ada OTP yang aktif. Silakan minta kirim ulang.")
	}

	// Cek 3: Apakah kedaluwarsa?
	if time.Now().After(user.VerificationExpiresAt.Time) {
		utils.AuthLogger.Printf("OTP verification failed (expired): %s", userUUID)
		return utils.SendError(c, fiber.StatusForbidden, "Kode OTP telah kedaluwarsa. Silakan minta kirim ulang.")
	}

	// Cek 4: Apakah OTP cocok?
	if user.VerificationToken.String != input.OTP {
		utils.AuthLogger.Printf("OTP verification failed (wrong OTP): %s", userUUID)
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

	utils.AuthLogger.Printf("OTP verification successful: %s", userUUID)
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
			utils.AuthLogger.Printf("Resend OTP failed (429 Too Many Requests): %s", userUUID)
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

	utils.AuthLogger.Printf("Resent verification OTP successful: %s", userUUID)
	return utils.SendSuccess(c, fiber.StatusOK, "Kode OTP baru telah dikirim ke email Anda.", nil)
}

func Login(c *fiber.Ctx) error {
	input := c.Locals("request_body").(*dto.LoginRequest)

	// 1. Cari user berdasarkan email
	var user models.User
	err := database.DB.Preload("Roles").Preload("Profile").First(&user, "email = ?", input.Email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.AuthLogger.Printf("Login failed (user not found): %s", input.Email)
		return utils.SendError(c, fiber.StatusUnauthorized, "Invalid credentials")
	}
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Database query error")
	}

	// 2. Cari provider auth 'local' untuk user ini
	var authProvider models.UserAuthProvider
	err = database.DB.Where("user_id = ? AND provider = ?", user.ID, "local").First(&authProvider).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.AuthLogger.Printf("Login failed (local auth not found): %s", input.Email)
		return utils.SendError(c, fiber.StatusUnauthorized, "Akun ini terdaftar melalui metode lain (misal: Google). Silakan login dengan metode tersebut.")
	}
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Database auth query error")
	}

	// 3. Verifikasi password
	match, err := utils.CheckPasswordHash(input.Password, authProvider.Password.String)
	if err != nil || !match {
		utils.AuthLogger.Printf("Login failed (invalid credentials): %s", input.Email)
		return utils.SendError(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	// 4. Generate JWT
	token, err := utils.GenerateJWT(user)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to generate JWT token")
	}

	go utils.AddUserToFrequentLoginFilter(user)
	responseData := dto.UserResponseJson(user, token)

	utils.AuthLogger.Printf("User login successful: %s (UUID: %s)", responseData.Email, responseData.ID)
	return utils.SendSuccess(c, fiber.StatusOK, "Login successful", responseData)
}

func LoginWithGoogle(c *fiber.Ctx) error {
	input := c.Locals("request_body").(*dto.GoogleLoginRequest)

	// 1. Verifikasi ID Token dengan Google
	userInfo, err := utils.VerifyGoogleIDToken(input.IDToken)
	if err != nil {
		utils.AuthLogger.Printf("Google login failed (token verification error): %v", err)
		return utils.SendError(c, fiber.StatusUnauthorized, err.Error())
	}

	googleEmail := userInfo.Email
	googleID := userInfo.Sub
	googleName := userInfo.Name

	if googleEmail == "" || googleID == "" {
		return utils.SendError(c, fiber.StatusUnauthorized, "Token Google tidak valid (data hilang)")
	}

	var user models.User
	var authProvider models.UserAuthProvider

	// Mulai Transaksi
	tx := database.DB.Begin()
	if tx.Error != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memulai transaksi")
	}
	defer tx.Rollback() // Rollback jika ada panic atau error

	// 2. Cari auth provider berdasarkan Google ID
	err = tx.Where("provider = ? AND provider_id = ?", "google", googleID).First(&authProvider).Error

	if err == nil {
		// --- KASUS 1: User ditemukan (Sudah pernah login via Google) ---
		if err = tx.Preload("Roles").Preload("Profile").First(&user, authProvider.UserID).Error; err != nil {
			return utils.SendError(c, fiber.StatusNotFound, "Data user terkait auth provider tidak ditemukan")
		}
		// Opsional: Update nama jika berubah di Google
		if user.Name != googleName {
			tx.Model(&user).Update("name", googleName)
		}

	} else if errors.Is(err, gorm.ErrRecordNotFound) {

		// 3. Cari user berdasarkan email (mungkin sudah punya akun 'local')
		err = tx.Preload("Roles").Preload("Profile").Where("email = ?", googleEmail).First(&user).Error

		if err == nil {
			// --- KASUS 2: User ditemukan berdasarkan Email (Akun 'local' ada) ---
			// Tautkan akun ini dengan Google ID
			authProvider = models.UserAuthProvider{
				UserID:     user.ID,
				Provider:   "google",
				ProviderID: googleID,
				Password:   sql.NullString{Valid: false},
			}
			if err := tx.Create(&authProvider).Error; err != nil {
				return utils.SendError(c, fiber.StatusInternalServerError, "Gagal menautkan akun Google")
			}
			// Update nama dan status verifikasi email
			tx.Model(&user).Updates(map[string]interface{}{
				"name":              googleName,
				"email_verified_at": sql.NullTime{Time: time.Now(), Valid: true},
			})

		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// --- KASUS 3: User tidak ditemukan (Akun baru) ---
			defaultRole, err := utils.GetRoleByName(string(constants.UserRole))
			if err != nil {
				return utils.SendError(c, fiber.StatusInternalServerError, "Konfigurasi server error")
			}

			// 3a. Buat User baru
			user = models.User{
				Name:            googleName,
				Email:           googleEmail,
				EmailVerifiedAt: sql.NullTime{Time: time.Now(), Valid: true}, // Email dari Google sudah terverifikasi
			}
			if err := tx.Create(&user).Error; err != nil {
				return utils.SendError(c, fiber.StatusInternalServerError, "Gagal membuat akun user baru")
			}

			// 3b. Buat Auth Provider untuk Google
			authProvider = models.UserAuthProvider{
				UserID:     user.ID,
				Provider:   "google",
				ProviderID: googleID,
				Password:   sql.NullString{Valid: false},
			}
			if err := tx.Create(&authProvider).Error; err != nil {
				return utils.SendError(c, fiber.StatusInternalServerError, "Gagal membuat auth provider")
			}

			// 3c. Tetapkan Role
			if err := tx.Model(&user).Association("Roles").Append(&defaultRole); err != nil {
				return utils.SendError(c, fiber.StatusInternalServerError, "Gagal menetapkan role")
			}
			user.Roles = []*models.Role{&defaultRole}
			user.Profile = models.Profile{}

			utils.AddEmailToRegistrationFilter(user.Email)

		} else {
			// Error database lain saat mencari by email
			return utils.SendError(c, fiber.StatusInternalServerError, "Database error (email lookup)")
		}

	} else {
		// Error database lain saat mencari by provider
		return utils.SendError(c, fiber.StatusInternalServerError, "Database error (provider lookup)")
	}

	// 5. Commit Transaksi
	if err := tx.Commit().Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal menyimpan data")
	}

	// 6. Generate JWT
	token, err := utils.GenerateJWT(user)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to generate JWT token")
	}

	responseData := dto.UserResponseJson(user, token)
	utils.AuthLogger.Printf("User login (Google) successful: %s (UUID: %s)", responseData.Email, responseData.ID)
	return utils.SendSuccess(c, fiber.StatusOK, "Login Google sukses", responseData)
}

func Logout(c *fiber.Ctx) error {
	userUUID, _ := c.Locals("user_id").(string)
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
		utils.AuthLogger.Printf("Logout skipped (token already expired): %s", userUUID)
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

	utils.AuthLogger.Printf("Logout successful: %s", userUUID)
	return utils.SendSuccess(c, fiber.StatusOK, "Logout successful", nil)
}
