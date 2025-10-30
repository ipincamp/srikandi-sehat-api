package handlers

import (
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"

	"github.com/gofiber/fiber/v2"
)

func GetNotificationHistory(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)

	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	var notifications []models.Notification
	database.DB.Where("user_id = ?", user.ID).Order("created_at desc").Find(&notifications)

	return utils.SendSuccess(c, fiber.StatusOK, "Notification history fetched", notifications)
}

func MarkNotificationAsRead(c *fiber.Ctx) error {
	// Ambil ID notifikasi dari parameter URL
	notificationID, err := c.ParamsInt("id")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid notification ID")
	}

	// Ambil UUID pengguna dari token JWT (disimpan oleh middleware)
	userUUID := c.Locals("user_id").(string)

	// Cari user berdasarkan UUID untuk mendapatkan ID integer-nya
	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	// Cari notifikasi berdasarkan ID dan UserID untuk memastikan pengguna hanya bisa mengubah notifikasinya sendiri
	var notification models.Notification
	result := database.DB.First(&notification, "id = ? AND user_id = ?", notificationID, user.ID)
	if result.Error != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Notification not found")
	}

	// Jika notifikasi belum dibaca, update statusnya
	if !notification.IsRead {
		database.DB.Model(&notification).Update("is_read", true)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Notification marked as read", nil)
}

// SendTestNotification sends a test FCM notification to the authenticated user.
// !! Should only be enabled in development environments. !!
func SendTestNotification(c *fiber.Ctx) error {
	// --- Environment Check (Simple) ---
	// Hanya izinkan jika APP_ENV tidak 'production'
	if config.Get("APP_ENV") == "production" {
		return utils.SendError(c, fiber.StatusForbidden, "This endpoint is disabled in production environment.")
	}
	// --- End Environment Check ---

	userUUID := c.Locals("user_id").(string)
	input := c.Locals("request_body").(*dto.TestNotificationRequest)

	var user models.User
	if err := database.DB.Select("id", "fcm_token").First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	if user.FcmToken == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "User does not have an FCM token registered.")
	}

	err := utils.SendFCMNotification(user.ID, user.FcmToken, input.Title, input.Body, input.Data)
	if err != nil {
		// Log error server-side
		utils.ErrorLogger.Printf("Failed to send test FCM notification to user %d (%s): %v", user.ID, userUUID, err)
		// Berikan pesan error yang lebih umum ke client
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to send test notification. Check server logs.")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Test notification sent successfully.", fiber.Map{
		"to_token": user.FcmToken,
		"title":    input.Title,
		"body":     input.Body,
		"data":     input.Data,
	})
}
