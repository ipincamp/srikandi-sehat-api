package handlers

import (
	"ipincamp/srikandi-sehat/database"
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
