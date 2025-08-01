package handlers

import (
	"errors"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"
	"log"

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

	if errors := utils.ValidateStruct(input); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	var existingUser models.User
	err := database.DB.First(&existingUser, "email = ?", input.Email).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.SendError(c, fiber.StatusConflict, "Email has already been taken")
	}

	hashedPassword, _ := utils.HashPassword(input.Password)

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return utils.SendError(c, fiber.StatusConflict, "Failed to create user")
	}

	var defaultRole models.Role
	if err := tx.First(&defaultRole, "name = ?", string(constants.UserRole)).Error; err != nil {
		tx.Rollback()
		log.Printf("CRITICAL: Default role '%s' not found. Please run the seeder.", constants.UserRole)
		return utils.SendError(c, fiber.StatusInternalServerError, "Server configuration error: default role not found")
	}

	if err := tx.Model(&user).Association("Roles").Append(&defaultRole); err != nil {
		tx.Rollback()
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to assign role to user")
	}

	if err := tx.Commit().Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	database.DB.Preload("Roles").First(&user, user.ID)

	responseData := dto.AuthResponseJson(user)

	return utils.SendSuccess(c, fiber.StatusCreated, "User registered successfully", responseData)
}

func Login(c *fiber.Ctx) error {
	input := new(dto.LoginRequest)
	if err := c.BodyParser(input); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if errors := utils.ValidateStruct(input); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	var user models.User
	result := database.DB.Preload("Roles.Permissions").Preload("Permissions").First(&user, "email = ?", input.Email)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) || !utils.CheckPasswordHash(input.Password, user.Password) {
		return utils.SendError(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Could not generate token")
	}

	responseData := dto.AuthResponseJson(user, token)

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
