package handlers

import (
	"errors"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Profile(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)
	var user models.User
	result := database.DB.Preload("Roles.Permissions").Preload("Permissions").First(&user, "uuid = ?", userUUID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	responseData := dto.AuthResponseJson(user)

	return utils.SendSuccess(c, fiber.StatusOK, "Profile fetched successfully", responseData)
}

func UpdateDetails(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)

	input := new(dto.UpdateDetailsRequest)
	if err := utils.ParseAndValidate(c, input); err != nil {
		return err
	}

	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	if err := database.DB.Model(&user).Updates(input).Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update user details")
	}

	responseData := dto.AuthResponseJson(user)
	return utils.SendSuccess(c, fiber.StatusOK, "Profile details updated successfully", responseData)
}

func ChangePassword(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)

	input := new(dto.ChangePasswordRequest)
	if err := utils.ParseAndValidate(c, input); err != nil {
		return err
	}

	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	if !utils.CheckPasswordHash(input.OldPassword, user.Password) {
		return utils.SendError(c, fiber.StatusBadRequest, "Old password does not match")
	}

	if input.NewPassword != input.ConfirmPassword {
		return utils.SendError(c, fiber.StatusBadRequest, "New password and confirmation do not match")
	}

	newHashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to process new password")
	}

	if err := database.DB.Model(&user).Update("password", newHashedPassword).Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update password")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Password changed successfully", nil)
}

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User

	adminUserIDs := database.DB.Table("user_roles").
		Select("user_roles.user_id").
		Joins("join roles on user_roles.role_id = roles.id").
		Where("roles.name = ?", string(constants.AdminRole))

	database.DB.Preload("Roles").
		Where("id NOT IN (?)", adminUserIDs).
		Find(&users)

	var responseData []dto.UserResponse
	for _, user := range users {
		responseData = append(responseData, dto.AuthResponseJson(user))
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Users fetched successfully", responseData)
}

func GetUserByID(c *fiber.Ctx) error {
	userUUID := c.Params("id")

	var user models.User
	result := database.DB.Preload("Roles").First(&user, "uuid = ?", userUUID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	responseData := dto.AuthResponseJson(user)

	return utils.SendSuccess(c, fiber.StatusOK, "User fetched successfully", responseData)
}
