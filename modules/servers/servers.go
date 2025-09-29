package servers

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"line-oa-backend/configs"
	"line-oa-backend/controllers"
	"line-oa-backend/db"
	"line-oa-backend/pkg/helper"
	"line-oa-backend/repositories/users"
	userService "line-oa-backend/service/users"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app    *fiber.App
	config *configs.Config
}

func NewServer() *Server {
	// Load configuration
	cfg := configs.Load()

	// Connect to database
	database, err := db.Connect(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	userRepo := users.NewUserRepository(database)

	// Initialize services
	userSvc := userService.NewUserService(userRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userSvc)

	// Initialize handlers
	handler := NewHandler(userController)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return helper.RespondWithError(c, code, err.Error())
		},
	})

	// Setup routes
	handler.SetupRoutes(app)

	return &Server{
		app:    app,
		config: cfg,
	}
}

func (s *Server) Start() {
	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		
		// Close database connection
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
		
		// Shutdown server
		if err := s.app.Shutdown(); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	log.Printf("Server starting on port %s", s.config.Port)

	if err := s.app.Listen(":" + s.config.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
