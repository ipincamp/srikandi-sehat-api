package middleware

import (
	"ipincamp/srikandi-sehat/src/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// MaintenanceMiddleware checks if the application is in maintenance mode.
func MaintenanceMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isMaintenance, message := utils.GetMaintenanceStatus()

		// Jika tidak dalam maintenance, selalu lanjutkan request
		if !isMaintenance {
			return c.Next()
		}

		// --- Maintenance sedang Aktif ---

		// 1. Selalu izinkan akses ke health check
		if c.Path() == "/api/health" {
			return c.Next()
		}

		// 2. Selalu izinkan request ke endpoint manajemen maintenance untuk diproses lebih lanjut.
		//    Middleware AuthMiddleware dan AdminMiddleware (yang didefinisikan di grup route)
		//    akan menangani otorisasi apakah user adalah admin atau bukan.
		if strings.HasPrefix(c.Path(), "/api/admin/maintenance") {
			return c.Next() // Lanjutkan ke middleware berikutnya (Auth, Admin)
		}

		// 3. Untuk SEMUA path LAINNYA, periksa apakah user terautentikasi dan ada di whitelist
		userUUID, userIsAuthenticated := c.Locals("user_id").(string)
		if userIsAuthenticated && utils.IsUserWhitelisted(userUUID) {
			return c.Next() // User terautentikasi dan whitelisted, boleh akses endpoint lain
		}

		// 4. Jika tidak memenuhi kondisi di atas (path bukan health/maintenance admin,
		//    atau user tidak terautentikasi/tidak whitelisted), blokir request.
		return utils.SendError(c, fiber.StatusServiceUnavailable, message)
	}
}
