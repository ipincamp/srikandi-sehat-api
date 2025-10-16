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
