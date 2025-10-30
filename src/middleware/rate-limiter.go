package middleware

import (
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func IPRateLimiter(max int, expiration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return utils.SendError(c, fiber.StatusTooManyRequests, "Too many requests, please try again later.")
		},
	})
}

func UserRateLimiter(max int, expiration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			userUUID, ok := c.Locals("user_id").(string)
			if !ok {
				return c.IP()
			}
			return userUUID
		},
		LimitReached: func(c *fiber.Ctx) error {
			return utils.SendError(c, fiber.StatusTooManyRequests, "Too many requests, please try again later.")
		},
	})
}

func LoginRateLimiter(max int, expiration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			input, ok := c.Locals("request_body").(*dto.LoginRequest)
			if !ok {
				return c.IP()
			}
			return input.Email
		},
		LimitReached: func(c *fiber.Ctx) error {
			return utils.SendError(c, fiber.StatusTooManyRequests, "Terlalu banyak percobaan login untuk akun ini. Silakan coba lagi nanti.")
		},
	})
}
