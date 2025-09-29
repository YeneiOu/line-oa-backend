package main

import (
	"log"

	"line-oa-backend/internal/adapters/handlers"
	"line-oa-backend/internal/adapters/middleware"
	"line-oa-backend/internal/adapters/repositories"
	"line-oa-backend/internal/adapters/services"
	"line-oa-backend/internal/application"
	"line-oa-backend/internal/infrastructure/config"
	"line-oa-backend/internal/infrastructure/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	userRepo := repositories.NewMongoUserRepository(db)
	bookingRepo := repositories.NewMongoBookingRepository(db)

	// Initialize external services
	lineOAuthService := services.NewLINEOAuthServiceAdapter(
		cfg.LINEChannelID,
		cfg.LINEChannelSecret,
		cfg.LINERedirectURI,
	)
	lineMessagingService := services.NewLINEMessagingServiceAdapter(cfg.LINEChannelAccessToken)
	jwtService := services.NewJWTServiceAdapter(cfg.JWTSecret)
	notificationService := services.NewNotificationServiceAdapter(lineMessagingService)

	// Initialize application services (use cases)
	authService := application.NewAuthService(userRepo, lineOAuthService, jwtService)
	bookingService := application.NewBookingService(bookingRepo, userRepo, notificationService)

	// Initialize HTTP handlers
	authHandler := handlers.NewAuthHandlerAdapter(authService)
	bookingHandler := handlers.NewBookingHandlerAdapter(bookingService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(middleware.CORSMiddleware(cfg.FrontendURL))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "LINE OA Backend is running",
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/callback", authHandler.Callback)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes
	protected := api.Group("/", middleware.AuthMiddleware(jwtService))

	// User routes
	protected.Get("/me", authHandler.GetUser)

	// Booking routes
	bookings := protected.Group("/bookings")
	bookings.Post("/", bookingHandler.CreateBooking)
	bookings.Get("/", bookingHandler.GetBookings)
	bookings.Get("/:id", bookingHandler.GetBooking)
	bookings.Put("/:id", bookingHandler.UpdateBooking)
	bookings.Delete("/:id", bookingHandler.CancelBooking)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Frontend URL: %s", cfg.FrontendURL)
	log.Printf("LINE Channel ID: %s", cfg.LINEChannelID)

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
