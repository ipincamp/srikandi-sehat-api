package utils

import "github.com/gofiber/fiber/v2"

type ResponseJson struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendSuccess(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(ResponseJson{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func SendError(c *fiber.Ctx, statusCode int, message string) error {
	if statusCode >= 500 {
	} else if statusCode >= 400 {
		InfoLogger.Printf("Client Error (Status %d): %s - Path: %s", statusCode, message, c.Path())
	}

	return c.Status(statusCode).JSON(ResponseJson{
		Status:  false,
		Message: message,
		Data:    nil,
	})
}
