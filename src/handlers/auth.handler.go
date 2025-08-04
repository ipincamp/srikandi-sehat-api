package handlers

import (
	"errors"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"
	"ipincamp/srikandi-sehat/src/workers"

	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {
	input := c.Locals("request_body").(*dto.RegisterRequest)

	if utils.CheckEmailExistsInRegistrationFilter(input.Email) {
		return utils.SendError(c, fiber.StatusUnprocessableEntity, "Email already registered")
	}

	job := workers.Job{
		RegistrationData: *input,
	}
	workers.JobQueue <- job

	return utils.SendSuccess(c, fiber.StatusAccepted, "Your account is being processed. Please try logging in after a while.", nil)
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

	invalidToken := models.InvalidToken{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}

	if err := database.DB.Create(&invalidToken).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			// Token is already blocklisted
		} else {
			return utils.SendError(c, fiber.StatusInternalServerError, "Failed to invalidate token")
		}
	}

	remainingDuration := time.Until(expiresAt)
	utils.AddToBlocklistCache(tokenString, remainingDuration)

	return utils.SendSuccess(c, fiber.StatusOK, "Logout successful", nil)
}
