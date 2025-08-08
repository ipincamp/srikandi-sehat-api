package middleware

import (
	"ipincamp/srikandi-sehat/src/utils"

	"github.com/gofiber/fiber/v2"
)

func RecoverMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				utils.LogPanic(r)

				if c.Response().StatusCode() == 0 {
					utils.SendError(c, fiber.StatusInternalServerError, "An unexpected error occurred. Please try again later.")
				}
			}
		}()

		return c.Next()
	}
}
