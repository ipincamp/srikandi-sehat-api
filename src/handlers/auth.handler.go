package handlers

import (
	"errors"
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

func Register(c *fiber.Ctx) error {
	input := c.Locals("request_body").(*dto.RegisterRequest)

	// 1. Cek email langsung ke Database (Source of Truth)
	var existingUser models.User
	err := database.DB.First(&existingUser, "email = ?", input.Email).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Email sudah ada di database.
		return utils.SendError(c, fiber.StatusConflict, "Email is already registered")
	}

	// 2. Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to hash password for %s: %v", input.Email, err)
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to process account")
	}

	// 3. Buat user baru
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

		defaultRole, err := utils.GetRoleByName(string(constants.UserRole))
		if err != nil {
			// Jika role 'user' tidak ada, ini adalah kesalahan server
			utils.ErrorLogger.Printf("FATAL: Default role '%s' not found in database", constants.UserRole)
			return err
		}

		if err := tx.Model(&user).Association("Roles").Append(&defaultRole); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		utils.ErrorLogger.Printf("Failed to create user %s in database transaction: %v", input.Email, err)
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to create account")
	}

	// 5. Tambahkan ke Bloom Filter (di memori)
	utils.AddEmailToRegistrationFilter(user.Email)

	// 6. Kembalikan respon sukses instan (201 Created)
	return utils.SendSuccess(c, fiber.StatusCreated, "Registration successful! Please log in.", nil)
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
