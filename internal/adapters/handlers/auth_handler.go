package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"line-oa-backend/internal/application"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AuthHandlerAdapter implements AuthHandler using Fiber
type AuthHandlerAdapter struct {
	authService *application.AuthService
}

// NewAuthHandlerAdapter creates a new auth handler adapter
func NewAuthHandlerAdapter(authService *application.AuthService) *AuthHandlerAdapter {
	return &AuthHandlerAdapter{
		authService: authService,
	}
}

type LoginResponse struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"`
}

type CallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type CallbackResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// Login generates LINE OAuth URL
func (h *AuthHandlerAdapter) Login(c *fiber.Ctx) error {
	// Generate a random state for security
	state, err := generateRandomState()
	if err != nil {
		log.Printf("Failed to generate state: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate authentication state",
		})
	}

	// Get LINE OAuth URL
	authURL := h.authService.GetAuthURL(state)

	return c.JSON(LoginResponse{
		AuthURL: authURL,
		State:   state,
	})
}

// Callback handles LINE OAuth callback
func (h *AuthHandlerAdapter) Callback(c *fiber.Ctx) error {
	var req CallbackRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Authorization code is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Authenticate with LINE
	user, token, err := h.authService.AuthenticateWithLINE(ctx, req.Code)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Authentication failed",
		})
	}

	return c.JSON(CallbackResponse{
		Token: token,
		User:  user,
	})
}

// RefreshToken refreshes the JWT token
func (h *AuthHandlerAdapter) RefreshToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Authorization header is required",
		})
	}

	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}

	tokenString := authHeader[7:] // Remove "Bearer " prefix

	newToken, err := h.authService.RefreshToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to refresh token",
		})
	}

	return c.JSON(fiber.Map{
		"token": newToken,
	})
}

// GetUser returns current user information
func (h *AuthHandlerAdapter) GetUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := h.authService.GetUser(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

func generateRandomState() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
