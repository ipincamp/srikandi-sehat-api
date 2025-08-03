package middleware

import (
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/utils"
	"slices"

	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware(c *fiber.Ctx) error {
	userUUID, ok := c.Locals("user_id").(string)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	roles, err := utils.GetUserRoles(userUUID)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	if slices.Contains(roles, string(constants.AdminRole)) {
		return c.Next()
	}

	return utils.SendError(c, fiber.StatusForbidden, "You do not have permission to access this resource")
}

func UserMiddleware(c *fiber.Ctx) error {
	userUUID, ok := c.Locals("user_id").(string)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	roles, err := utils.GetUserRoles(userUUID)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	if slices.Contains(roles, string(constants.UserRole)) {
		return c.Next()
	}

	return utils.SendError(c, fiber.StatusForbidden, "You do not have permission to access this resource")
}
