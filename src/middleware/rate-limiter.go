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
			var input dto.LoginRequest

			// Ambil body mentah
			body := c.Body()
			if len(body) == 0 {
				return c.IP() // Fallback ke IP jika body kosong
			}

			// Unmarshal JSON body ke struct
			if err := json.Unmarshal(body, &input); err != nil {
				return c.IP() // Fallback ke IP jika JSON tidak valid
			}

			if input.Email == "" {
				return c.IP() // Fallback ke IP jika email kosong
			}

			// Kunci limiter sekarang adalah alamat email
			return input.Email
		},
		LimitReached: func(c *fiber.Ctx) error {
			return utils.SendError(c, fiber.StatusTooManyRequests, "Terlalu banyak percobaan login untuk akun ini. Silakan coba lagi nanti.")
		},
	})
}
