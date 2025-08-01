package middleware

import (
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"

	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware(c *fiber.Ctx) error {
	userUUID, ok := c.Locals("user_id").(string)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	var user models.User
	if err := database.DB.Preload("Roles").First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	for _, role := range user.Roles {
		if role.Name == string(constants.AdminRole) {
			return c.Next()
		}
	}

	return utils.SendError(c, fiber.StatusForbidden, "Administrator access required")
}
