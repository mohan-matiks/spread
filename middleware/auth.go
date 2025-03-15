package middleware

import (
	"strings"

	"github.com/SwishHQ/spread/config"
	"github.com/SwishHQ/spread/logger"
	"github.com/SwishHQ/spread/src/service"
	"github.com/SwishHQ/spread/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func AuthMiddleware(c *fiber.Ctx, userService service.UserService) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		logger.L.Error("In authMiddleware: Missing Authorization header")
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}

	// Extract the token (Bearer <token>)
	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		logger.L.Error("In authMiddleware: Invalid Authorization format")
		return utils.UnauthorizedResponse(c, "Invalid Authorization format")
	}

	// Parse JWT
	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(config.TokenSecret), nil
	})
	if err != nil || !token.Valid {
		logger.L.Error("In authMiddleware: Invalid token", zap.Error(err))
		return utils.UnauthorizedResponse(c, "Invalid token")
	}

	// Extract user ID from token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.L.Error("In authMiddleware: Invalid token claims")
		return utils.UnauthorizedResponse(c, "Invalid token claims")
	}

	userID, ok := claims["id"].(string)
	if !ok {
		logger.L.Error("In authMiddleware: Invalid user ID in token")
		return utils.UnauthorizedResponse(c, "Invalid user ID in token")
	}

	// Fetch user from database
	user, err := userService.GetUser(userID)
	if err != nil {
		logger.L.Error("In authMiddleware: Error getting user", zap.Error(err))
		return utils.UnauthorizedResponse(c, "User not found")
	}

	// Store user in context
	c.Locals("user", user)
	return c.Next()
}
