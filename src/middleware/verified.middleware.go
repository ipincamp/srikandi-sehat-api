package middleware

import (
	"database/sql"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"

	"github.com/gofiber/fiber/v2"
)

func VerifiedMiddleware(c *fiber.Ctx) error {
	userUUID, ok := c.Locals("user_id").(string)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	var user struct {
		EmailVerifiedAt sql.NullTime
	}

	// Hanya untuk status verifikasi
	// TODO: Ini bisa di-cache di utils/cache.util.go untuk performa lebih baik
	if err := database.DB.Model(&models.User{}).
		Select("email_verified_at").
		Where("uuid = ?", userUUID).
		First(&user).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	if !user.EmailVerifiedAt.Valid {
		return utils.SendError(c, fiber.StatusForbidden, "Please verify your email address to access this resource.")
	}

	return c.Next()
}
