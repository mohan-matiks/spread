package middleware

import (
	"github.com/SwishHQ/spread/logger"
	"github.com/SwishHQ/spread/src/service"
	"github.com/SwishHQ/spread/utils"
	"github.com/gofiber/fiber/v2"
)

func AuthKeyMiddleware(c *fiber.Ctx, authKeyService service.AuthKeyService) error {
	authKey := c.Get("x-auth-key")
	if authKey == "" {
		logger.L.Error("In AuthKeyMiddleware: No auth key found")
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}
	valid, err := authKeyService.ValidateAuthKey(authKey)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	if !valid {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}
	return c.Next()
}
