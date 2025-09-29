package helper

import (
	"line-oa-backend/entities"

	"github.com/gofiber/fiber/v2"
)

// RespondWithJSON sends a standardized JSON response
func RespondWithJSON(c *fiber.Ctx, status string, statusCode int, message string, data interface{}) error {
	response := entities.Responses{
		Status:  status,
		Code:    statusCode,
		Message: message,
		Data:    data,
	}
	return c.Status(statusCode).JSON(response)
}

// RespondWithError sends a standardized error response
func RespondWithError(c *fiber.Ctx, statusCode int, message string) error {
	return RespondWithJSON(c, "error", statusCode, message, nil)
}

// RespondWithSuccess sends a standardized success response
func RespondWithSuccess(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return RespondWithJSON(c, "success", statusCode, message, data)
}
