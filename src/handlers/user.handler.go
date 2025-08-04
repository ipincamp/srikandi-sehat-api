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

	if input.Name != nil {
		if err := tx.Model(&user).Update("name", *input.Name).Error; err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update user name")
		}
	}

	updateData := make(map[string]interface{})
	if input.PhoneNumber != nil {
		updateData["phone_number"] = *input.PhoneNumber
	}
	if input.HeightCM != nil {
		updateData["height_cm"] = *input.HeightCM
	}
	if input.WeightKG != nil {
		updateData["weight_kg"] = *input.WeightKG
	}
	if input.LastEducation != nil {
		updateData["last_education"] = *input.LastEducation
	}
	if input.ParentLastEducation != nil {
		updateData["parent_last_education"] = *input.ParentLastEducation
	}
	if input.ParentLastJob != nil {
		updateData["parent_last_job"] = *input.ParentLastJob
	}
	if input.InternetAccess != nil {
		updateData["internet_access"] = *input.InternetAccess
	}
	if input.MenarcheAge != nil {
		updateData["menarche_age"] = *input.MenarcheAge
	}

	if input.DateOfBirth != nil {
		if dob, err := time.Parse("2006-01-02", *input.DateOfBirth); err == nil {
			updateData["date_of_birth"] = &dob
		}
	}
	if input.VillageCode != nil {
		var village region.Village
		if err := tx.First(&village, "code = ?", *input.VillageCode).Error; err != nil {
			return utils.SendError(c, fiber.StatusNotFound, "Village code not found")
		}
		updateData["village_id"] = &village.ID
	}

	var profile models.Profile
	err := tx.Where(models.Profile{UserID: user.ID}).First(&profile).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		updateData["user_id"] = user.ID
		if err := tx.Model(&models.Profile{}).Create(updateData).Error; err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Failed to create profile")
		}
	} else if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to find profile")
	} else {
		if len(updateData) > 0 {
			if err := tx.Model(&profile).Updates(updateData).Error; err != nil {
				return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update profile")
			}
		}
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
	queries := c.Locals("request_queries").(*dto.UserQuery)

	var users []models.User

	adminUserIDs := database.DB.Table("user_roles").
		Select("user_roles.user_id").
		Joins("join roles on user_roles.role_id = roles.id").
		Where("roles.name = ?", string(constants.AdminRole))

	query := database.DB.Model(&models.User{}).
		Joins("JOIN profiles ON users.id = profiles.user_id").
		Joins("JOIN villages ON profiles.village_id = villages.id").
		Joins("JOIN classifications ON villages.classification_id = classifications.id")

	if queries.Classification != "" {
		query = query.Where("classifications.name = ?", queries.Classification)
	}

	query = query.Where("users.id NOT IN (?)", adminUserIDs)

	page := queries.Page
	if page == 0 {
		page = 1
	}
	limit := queries.Limit
	if limit == 0 {
		limit = 10
	}

	pagination, paginateScope := utils.GeneratePagination(page, limit, query, &models.User{})

	query.Select("users.uuid, users.name, users.created_at").Scopes(paginateScope).Find(&users)

	var responseData []fiber.Map
	for _, user := range users {
		responseData = append(responseData, fiber.Map{
			"id":         user.UUID,
			"name":       user.Name,
			"created_at": user.CreatedAt,
		})
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
