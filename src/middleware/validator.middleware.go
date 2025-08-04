package middleware

import (
	"ipincamp/srikandi-sehat/src/utils"

	"github.com/gofiber/fiber/v2"
)

func ValidateBody[T any](c *fiber.Ctx) error {
	body := new(T)

	if err := c.BodyParser(body); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if validationErrors := utils.ValidateStruct(body); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Validation failed",
			"errors":  validationErrors,
		})
	}

	c.Locals("request_body", body)

	return c.Next()
}

func ValidateQuery[T any](c *fiber.Ctx) error {
	queries := new(T)

	if err := c.QueryParser(queries); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid query parameters")
	}

	if validationErrors := utils.ValidateStruct(queries); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Validation failed",
			"errors":  validationErrors,
		})
	}

	c.Locals("request_queries", queries)

	return c.Next()
}
