package handlers

import (
	"errors"
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/models"
	"ipincamp/srikandi-sehat/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterInput struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateJWT(user models.User) (string, error) {
	secretKey := config.Get("JWT_SECRET")

	// Take expiration time from config, default to 24 hours if not set
	expHoursStr := config.Get("JWT_EXPIRATION_HOURS")
	expHours, err := strconv.Atoi(expHoursStr)
	if err != nil {
		expHours = 24 // Default
	}

	claims := jwt.MapClaims{
		"uid": user.UUID,
		"exp": time.Now().Add(time.Hour * time.Duration(expHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func Register(c *fiber.Ctx) error {
	input := new(RegisterInput)
	if err := utils.ParseAndValidate(c, input); err != nil {
		return err
	}

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Could not hash password")
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		return utils.ErrorResponse(c, fiber.StatusConflict, "Email already exists")
	}

	userData := fiber.Map{
		"id":         user.UUID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "User registered successfully", userData)
}

func Login(c *fiber.Ctx) error {
	input := new(LoginInput)
	if err := utils.ParseAndValidate(c, input); err != nil {
		return err
	}

	var user models.User
	result := database.DB.First(&user, "email = ?", input.Email)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) || !checkPasswordHash(input.Password, user.Password) {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	token, err := generateJWT(user)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Could not generate token")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Login successful", fiber.Map{"token": token})
}

func Profile(c *fiber.Ctx) error {
	userUUID := c.Locals("uid").(string)
	var user models.User
	result := database.DB.Select("uuid", "name", "email", "created_at", "updated_at").First(&user, "uuid = ?", userUUID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Profile fetched successfully", user)
}

func Logout(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, "Successfully logged out", nil)
}
