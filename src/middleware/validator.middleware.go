package middleware

import (
	"ipincamp/srikandi-sehat/src/utils"

	"github.com/gofiber/fiber/v2"
)

func Validate[T any](c *fiber.Ctx) error {
	input := new(T)

	if err := c.BodyParser(input); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if validationErrors := utils.ValidateStruct(input); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Validation failed",
			"errors":  validationErrors,
		})
	}

	c.Locals("request", input)

	return c.Next()
}
