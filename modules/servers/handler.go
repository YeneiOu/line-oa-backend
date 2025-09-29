package servers

import (
	"line-oa-backend/controllers"
	"line-oa-backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Handler struct {
	userController *controllers.UserController
}

func NewHandler(userController *controllers.UserController) *Handler {
	return &Handler{
		userController: userController,
	}
}

func (h *Handler) SetupRoutes(app *fiber.App) {
	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(utils.CORSMiddleware())

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "LINE OA Backend is running",
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// User routes
	users := api.Group("/users")
	users.Post("/", h.userController.CreateUser)
	users.Get("/:id", h.userController.GetUser)
	
	// Profile route
	api.Get("/profile", h.userController.GetUserProfile)
}
