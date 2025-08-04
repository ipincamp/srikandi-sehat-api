package handlers

import (
	"errors"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/models/region"
	"ipincamp/srikandi-sehat/src/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetMyProfile(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)

	var user models.User

	err := database.DB.
		Preload("Roles").
		Preload("Profile.Village.Classification").
		Preload("Profile.Village.District.Regency.Province").
		First(&user, "uuid = ?", userUUID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SendError(c, fiber.StatusNotFound, "User not found")
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Database error")
	}

	if user.Profile.ID == 0 {
		return utils.SendError(c, fiber.StatusNotFound, "Your profile has not been created yet. Please update your profile first.")
	}

	responseData := dto.UserResponseJson(user)
	return utils.SendSuccess(c, fiber.StatusOK, "Profile fetched successfully", responseData)
}

func UpdateOrCreateProfile(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)

	input := c.Locals("request_body").(*dto.UpdateProfileRequest)

	tx := database.DB.Begin()
	if tx.Error != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}
	defer tx.Rollback()

	var user models.User
	if err := tx.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	if input.Name != "" {
		if err := tx.Model(&user).Update("name", input.Name).Error; err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update user name")
		}
	}

	dob, _ := time.Parse("2006-01-02", input.DateOfBirth)

	var village region.Village
	if err := tx.First(&village, "code = ?", input.VillageCode).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Village code not found")
	}

	profileData := models.Profile{
		UserID:              user.ID,
		PhoneNumber:         input.PhoneNumber,
		DateOfBirth:         &dob,
		HeightCM:            input.HeightCM,
		WeightKG:            input.WeightKG,
		AddressStreet:       input.AddressStreet,
		VillageID:           &village.ID,
		LastEducation:       input.LastEducation,
		ParentLastEducation: input.ParentLastEducation,
		ParentLastJob:       input.ParentLastJob,
		InternetAccess:      input.InternetAccess,
		MenarcheAge:         input.MenarcheAge,
	}

	var existingProfile models.Profile
	err := tx.Where(models.Profile{UserID: user.ID}).
		Assign(profileData).
		FirstOrCreate(&existingProfile).Error

	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to save profile")
	}

	if err := tx.Commit().Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Profile updated successfully", nil)
}

func ChangeMyPassword(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)

	input := c.Locals("request_body").(*dto.ChangePasswordRequest)

	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SendError(c, fiber.StatusNotFound, "User not found")
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Database error")
	}

	match, err := utils.CheckPasswordHash(input.OldPassword, user.Password)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to verify old password")
	}
	if !match {
		return utils.SendError(c, fiber.StatusUnauthorized, "Old password is incorrect")
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

	query := database.DB.Preload("Roles").Where("id NOT IN (?)", adminUserIDs)

	pagination, paginateScope := utils.GeneratePagination(c, query, &models.User{})

	query.Scopes(paginateScope).Find(&users)

	var responseData []dto.UserResponse
	for _, user := range users {
		responseData = append(responseData, dto.UserResponseJson(user))
	}

	paginatedResponse := fiber.Map{
		"data": responseData,
		"meta": pagination,
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Users fetched successfully", paginatedResponse)
}

func GetUserByID(c *fiber.Ctx) error {
	userUUID := c.Params("id")

	var user models.User
	result := database.DB.Preload("Roles").First(&user, "uuid = ?", userUUID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	responseData := dto.UserResponseJson(user)

	return utils.SendSuccess(c, fiber.StatusOK, "User fetched successfully", responseData)
}
