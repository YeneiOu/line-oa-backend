package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORSMiddleware creates CORS middleware
func CORSMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		// AllowOrigins:     frontendURL + ",http://localhost:3000,http://127.0.0.1:3000",
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	})
}
