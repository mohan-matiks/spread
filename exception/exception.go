package exception

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ResourceNotFoundException(resourceName string, fieldName string, fieldValue string) error {
	msg := fmt.Sprintf("%s not found with %s : %s", resourceName, fieldName, fieldValue)
	return fiber.NewError(fiber.StatusNotFound, msg)
}

func BadRequestException(msg string) error {
	return fiber.NewError(fiber.StatusBadRequest, msg)
}

func ConflictException(resourceName string, fieldName string, fieldValue string) error {
	msg := fmt.Sprintf("%s with %s : %s already exists", resourceName, fieldName, fieldValue)
	return fiber.NewError(fiber.StatusConflict, msg)
}

func UnauthorizedException() error {
	return fiber.ErrUnauthorized
}

func InternalServerErrorException() error {
	msg := fmt.Sprintf("Something went wrong. Please try again later.")
	return fiber.NewError(fiber.StatusInternalServerError, msg)
}
