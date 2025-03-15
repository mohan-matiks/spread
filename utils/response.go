package utils

import "github.com/gofiber/fiber/v2"

func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	if data == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func ErrorResponse(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}

func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}

func ValidationErrorResponse(c *fiber.Ctx, errors []string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"errors":  errors,
	})
}
