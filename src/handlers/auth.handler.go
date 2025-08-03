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
	input := new(dto.RegisterRequest)
	if err := c.BodyParser(input); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}
	if validationErrors := utils.ValidateStruct(input); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Validation failed",
			"errors":  validationErrors,
		})
	}

	job := workers.Job{
		RegistrationData: *input,
	}
	workers.JobQueue <- job

	return utils.SendSuccess(c, fiber.StatusAccepted, "Your account is being processed. Please try logging in after a while.", nil)
}

func Login(c *fiber.Ctx) error {
	input := new(dto.LoginRequest)
	if err := c.BodyParser(input); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}
	if validationErrors := utils.ValidateStruct(input); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Validation failed",
			"errors":  validationErrors,
		})
	}

	var user models.User
	result := database.DB.Preload("Roles").Preload("Profile").First(&user, "email = ?", input.Email)

	match, err := utils.CheckPasswordHash(input.Password, user.Password)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || err != nil || !match {
		return utils.SendError(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Could not generate token")
	}

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
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to parse token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return utils.SendError(c, fiber.StatusInternalServerError, "Invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return utils.SendError(c, fiber.StatusInternalServerError, "Invalid expiration time in token")
	}
	expiresAt := time.Unix(int64(exp), 0)

	invalidToken := models.InvalidToken{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}

	if err := database.DB.Create(&invalidToken).Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to logout")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Logout successful", nil)
}
