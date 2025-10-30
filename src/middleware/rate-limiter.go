package middleware

import (
	"encoding/json"
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
			utils.ErrorLogger.Printf("Rate limit (IP) reached for IP: %s - Path: %s", c.IP(), c.Path())
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
			userUUID, _ := c.Locals("user_id").(string)
			utils.ErrorLogger.Printf("Rate limit (User) reached for User: %s (IP: %s) - Path: %s", userUUID, c.IP(), c.Path())
			return utils.SendError(c, fiber.StatusTooManyRequests, "Too many requests, please try again later.")
		},
	})
}

func LoginRateLimiter(max int, expiration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			var input dto.LoginRequest

			body := c.Body()
			if len(body) == 0 {
				return c.IP()
			}

			if err := json.Unmarshal(body, &input); err != nil {
				return c.IP()
			}

			if input.Email == "" {
				return c.IP()
			}

			return input.Email
		},
		LimitReached: func(c *fiber.Ctx) error {
			utils.ErrorLogger.Printf("Rate limit (Login) reached for IP: %s - Path: %s", c.IP(), c.Path())
			return utils.SendError(c, fiber.StatusTooManyRequests, "Terlalu banyak percobaan login untuk akun ini. Silakan coba lagi nanti.")
		},
	})
}
