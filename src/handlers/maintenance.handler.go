package handlers

import (
	"errors"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const maintenanceActiveKey = "maintenance_mode_active"
const maintenanceMessageKey = "maintenance_message"

// GetMaintenanceStatus retrieves the current maintenance status.
func GetMaintenanceStatus(c *fiber.Ctx) error {
	status, message := utils.GetMaintenanceStatus()
	return utils.SendSuccess(c, fiber.StatusOK, "Maintenance status fetched", dto.MaintenanceStatusResponse{
		IsMaintenance: status,
		Message:       message,
	})
}

// ToggleMaintenanceMode enables or disables maintenance mode. (Admin only)
func ToggleMaintenanceMode(c *fiber.Ctx) error {
	input := c.Locals("request_body").(*dto.ToggleMaintenanceRequest)

	tx := database.DB.Begin()
	if tx.Error != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}
	defer tx.Rollback() // Pastikan rollback jika terjadi error atau panic

	// --- Update active status menggunakan Save ---
	activeStr := strconv.FormatBool(*input.Active)
	settingActive := models.Setting{Key: maintenanceActiveKey, Value: activeStr}
	// Ganti .Model().Where().Update() dengan .Save()
	if err := tx.Save(&settingActive).Error; err != nil {
		utils.ErrorLogger.Printf("Error updating maintenance status: %v", err) // Tambahkan log error
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update maintenance status")
	}

	// --- Update message menggunakan Save (ini sudah benar sebelumnya, tapi kita konsistenkan) ---
	message := "Server is currently under maintenance. Please try again later." // Default message
	if input.Message != "" {
		message = input.Message
	}
	settingMessage := models.Setting{Key: maintenanceMessageKey, Value: message}
	if err := tx.Save(&settingMessage).Error; err != nil { // Tetap gunakan Save
		utils.ErrorLogger.Printf("Error updating maintenance message: %v", err) // Tambahkan log error
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update maintenance message")
	}

	if err := tx.Commit().Error; err != nil {
		utils.ErrorLogger.Printf("Error committing maintenance toggle transaction: %v", err) // Log error commit
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	// Reload cache setelah commit berhasil
	utils.ReloadMaintenanceStatus()

	statusMsg := "disabled"
	if *input.Active {
		statusMsg = "enabled"
	}

	// Ambil status terbaru dari cache untuk konsistensi response
	currentStatus, currentMessage := utils.GetMaintenanceStatus()

	return utils.SendSuccess(c, fiber.StatusOK, "Maintenance mode "+statusMsg+" successfully", dto.MaintenanceStatusResponse{
		IsMaintenance: currentStatus,  // Gunakan status dari cache yang baru direload
		Message:       currentMessage, // Gunakan pesan dari cache yang baru direload
	})
}

// AddUserToWhitelist adds a user to the maintenance bypass list. (Admin only)
func AddUserToWhitelist(c *fiber.Ctx) error {
	input := c.Locals("request_body").(*dto.WhitelistUserRequest)

	var user models.User
	if err := database.DB.Select("id").First(&user, "uuid = ?", input.UserUUID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SendError(c, fiber.StatusNotFound, "User with the specified UUID not found")
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Database error")
	}

	whitelistEntry := models.MaintenanceWhitelist{UserID: user.ID}
	result := database.DB.FirstOrCreate(&whitelistEntry)
	if result.Error != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to add user to whitelist")
	}

	if result.RowsAffected == 0 {
		return utils.SendSuccess(c, fiber.StatusOK, "User is already in the whitelist", nil)
	}

	// Reload whitelist cache
	utils.ReloadMaintenanceWhitelist()

	return utils.SendSuccess(c, fiber.StatusCreated, "User added to maintenance whitelist successfully", nil)
}

// RemoveUserFromWhitelist removes a user from the maintenance bypass list. (Admin only)
func RemoveUserFromWhitelist(c *fiber.Ctx) error {
	input := c.Locals("request_body").(*dto.WhitelistUserRequest)

	var user models.User
	if err := database.DB.Select("id").First(&user, "uuid = ?", input.UserUUID).Error; err != nil {
		// Even if user doesn't exist, try removing potential stale whitelist entry
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorLogger.Printf("Error finding user to remove from whitelist: %v", err)
		}
	}

	// Attempt to delete regardless of whether user was found, to handle stale entries
	result := database.DB.Where("user_id = (SELECT id FROM users WHERE uuid = ?)", input.UserUUID).Delete(&models.MaintenanceWhitelist{})

	if result.Error != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to remove user from whitelist")
	}

	if result.RowsAffected == 0 {
		return utils.SendError(c, fiber.StatusNotFound, "User not found in the whitelist")
	}

	// Reload whitelist cache
	utils.ReloadMaintenanceWhitelist()

	return utils.SendSuccess(c, fiber.StatusOK, "User removed from maintenance whitelist successfully", nil)
}

// GetWhitelistedUsers retrieves the list of users allowed during maintenance. (Admin only)
func GetWhitelistedUsers(c *fiber.Ctx) error {
	var whitelistedEntries []models.MaintenanceWhitelist
	if err := database.DB.Preload("User").Find(&whitelistedEntries).Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve whitelist")
	}

	response := make([]dto.WhitelistedUserResponse, len(whitelistedEntries))
	for i, entry := range whitelistedEntries {
		response[i] = dto.WhitelistedUserResponse{
			UserUUID: entry.User.UUID,
			UserName: entry.User.Name,
			Email:    entry.User.Email,
		}
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Whitelisted users fetched successfully", response)
}

// HealthCheck checks the application status including maintenance mode.
func HealthCheck(c *fiber.Ctx) error {
	isMaintenance, _ := utils.GetMaintenanceStatus()
	status := "OK"
	statusCode := http.StatusOK
	if isMaintenance {
		status = "MAINTENANCE"
		statusCode = http.StatusServiceUnavailable
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"status": status,
	})
}
