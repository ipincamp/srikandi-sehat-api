package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type ValidationError struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

func ValidateStruct(payload interface{}) []*ValidationError {
	var errors []*ValidationError
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func ParseAndValidate(c *fiber.Ctx, payload interface{}) error {
	if err := c.BodyParser(payload); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if errors := ValidateStruct(payload); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Validation failed",
			"data":    errors,
		})
	}
	return nil
}
