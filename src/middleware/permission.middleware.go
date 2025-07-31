package middleware

import (
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"

	"github.com/gofiber/fiber/v2"
)

func UserHasPermission(userUUID string, requiredPermission string) bool {
	var user models.User
	err := database.DB.Preload("Roles.Permissions").Preload("Permissions").First(&user, "uuid = ?", userUUID).Error
	if err != nil {
		return false
	}

	for _, p := range user.Permissions {
		if p.Name == requiredPermission {
			return true
		}
	}

	for _, r := range user.Roles {
		for _, p := range r.Permissions {
			if p.Name == requiredPermission {
				return true
			}
		}
	}

	return false
}

func PermissionMiddleware(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userUUID, ok := c.Locals("user_id").(string)
		if !ok {
			return utils.SendError(c, fiber.StatusUnauthorized, "Unauthorized")
		}

		if !UserHasPermission(userUUID, requiredPermission) {
			return utils.SendError(c, fiber.StatusForbidden, "You don't have the required permission")
		}

		return c.Next()
	}
}
