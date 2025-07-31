package handlers

import (
	"errors"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {
	input := new(dto.Register)
	if err := utils.ParseAndValidate(c, input); err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Could not hash password")
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		return utils.SendError(c, fiber.StatusConflict, "Email already exists")
	}

	userData := fiber.Map{
		"id":         user.UUID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	}

	return utils.SendSuccess(c, fiber.StatusCreated, "User registered successfully", userData)
}

func Login(c *fiber.Ctx) error {
	input := new(dto.Login)
	if err := utils.ParseAndValidate(c, input); err != nil {
		return err
	}

	var user models.User
	result := database.DB.First(&user, "email = ?", input.Email)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) || !utils.CheckPasswordHash(input.Password, user.Password) {
		return utils.SendError(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Could not generate token")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Login successful", fiber.Map{"token": token})
}

func Profile(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)
	var user models.User
	result := database.DB.Select("uuid", "name", "email", "created_at", "updated_at").First(&user, "uuid = ?", userUUID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Profile fetched successfully", user)
}

func Logout(c *fiber.Ctx) error {
	return utils.SendSuccess(c, fiber.StatusOK, "Successfully logged out", nil)
}
