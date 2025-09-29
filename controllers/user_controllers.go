package controllers

import (
	"context"
	"line-oa-backend/pkg/helper"
	"line-oa-backend/service/users"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService *users.UserService
}

func NewUserController(userService *users.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// CreateUser creates a new user
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	var req struct {
		LineUserID string `json:"line_user_id"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		PictureURL string `json:"picture_url"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return helper.RespondWithError(ctx, fiber.StatusBadRequest, "Invalid request body")
	}

	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := c.userService.CreateUser(context, req.LineUserID, req.Name, req.Email, req.PictureURL)
	if err != nil {
		return helper.RespondWithError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return helper.RespondWithSuccess(ctx, fiber.StatusCreated, "User created successfully", user)
}

// GetUser gets a user by ID
func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")
	if userID == "" {
		return helper.RespondWithError(ctx, fiber.StatusBadRequest, "User ID is required")
	}

	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := c.userService.GetUser(context, userID)
	if err != nil {
		return helper.RespondWithError(ctx, fiber.StatusNotFound, "User not found")
	}

	return helper.RespondWithSuccess(ctx, fiber.StatusOK, "User retrieved successfully", user)
}

// GetUserProfile gets current user profile (placeholder for auth)
func (c *UserController) GetUserProfile(ctx *fiber.Ctx) error {
	// This would normally get user ID from JWT token
	// For now, return a simple response
	return helper.RespondWithSuccess(ctx, fiber.StatusOK, "User profile endpoint", nil)
}
