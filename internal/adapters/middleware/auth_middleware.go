package middleware

import (
	"line-oa-backend/internal/ports"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware creates authentication middleware
func AuthMiddleware(jwtService ports.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		// Extract token
		tokenString := authHeader[7:] // Remove "Bearer " prefix

		// Validate token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Store user information in context
		c.Locals("userID", claims.UserID)
		c.Locals("lineUserID", claims.LineUserID)

		return c.Next()
	}
}
