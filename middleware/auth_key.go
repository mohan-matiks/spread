package middleware

import (
	"github.com/SwishHQ/spread/logger"
	"github.com/SwishHQ/spread/src/service"
	"github.com/SwishHQ/spread/utils"
	"github.com/gofiber/fiber/v2"
)

func AuthKeyMiddleware(c *fiber.Ctx, authKeyService service.AuthKeyService) error {
	key := c.Get("x-auth-key")
	if key == "" {
		logger.L.Error("In AuthKeyMiddleware: No auth key found")
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}
	authKey, err := authKeyService.GetByAuthKey(key)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	if authKey == nil {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}
	c.Locals("authKey", authKey)
	return c.Next()
}
