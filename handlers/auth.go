package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"line-oa-backend/database"
	"line-oa-backend/models"
	"line-oa-backend/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	lineOAuth  *services.LINEOAuthService
	jwtService *services.JWTService
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
	User  models.User `json:"user"`
}

func NewAuthHandler(lineOAuth *services.LINEOAuthService, jwtService *services.JWTService) *AuthHandler {
	return &AuthHandler{
		lineOAuth:  lineOAuth,
		jwtService: jwtService,
	}
}

// Login generates LINE OAuth URL
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Generate a random state for security
	state, err := generateRandomState()
	if err != nil {
		log.Printf("Failed to generate state: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate authentication state",
		})
	}

	// Get LINE OAuth URL
	authURL := h.lineOAuth.GetAuthURL(state)

	return c.JSON(LoginResponse{
		AuthURL: authURL,
		State:   state,
	})
}

// Callback handles LINE OAuth callback
func (h *AuthHandler) Callback(c *fiber.Ctx) error {
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

	// Exchange code for token
	tokenResp, err := h.lineOAuth.ExchangeCodeForToken(req.Code)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to exchange authorization code",
		})
	}

	// Get user profile
	profile, err := h.lineOAuth.GetProfile(tokenResp.AccessToken)
	if err != nil {
		log.Printf("Failed to get user profile: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to get user profile",
		})
	}

	// Find or create user in database
	db := database.GetDatabase()
	usersCollection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	// Try to find existing user
	err = usersCollection.FindOne(ctx, bson.M{"line_user_id": profile.UserID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User doesn't exist, create new one
			user = models.User{
				ID:         primitive.NewObjectID(),
				LineUserID: profile.UserID,
				Name:       profile.DisplayName,
				PictureURL: profile.PictureURL,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			_, err := usersCollection.InsertOne(ctx, user)
			if err != nil {
				log.Printf("Failed to create user: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create user account",
				})
			}
		} else {
			log.Printf("Failed to find user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}
	} else {
		// User exists, update profile information
		update := bson.M{
			"$set": bson.M{
				"name":        profile.DisplayName,
				"picture_url": profile.PictureURL,
				"updated_at":  time.Now(),
			},
		}

		_, err := usersCollection.UpdateOne(ctx, bson.M{"_id": user.ID}, update)
		if err != nil {
			log.Printf("Failed to update user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update user account",
			})
		}

		// Update local user object
		user.Name = profile.DisplayName
		user.PictureURL = profile.PictureURL
		user.UpdatedAt = time.Now()
	}

	// Generate JWT token
	jwtToken, err := h.jwtService.GenerateToken(user.ID, user.LineUserID)
	if err != nil {
		log.Printf("Failed to generate JWT token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate authentication token",
		})
	}

	return c.JSON(CallbackResponse{
		Token: jwtToken,
		User:  user,
	})
}

// Me returns current user information
func (h *AuthHandler) User(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID)

	db := database.GetDatabase()
	usersCollection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	return c.JSON(user)
}

// RefreshToken refreshes the JWT token
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Authorization header is required",
		})
	}

	tokenString := authHeader[7:] // Remove "Bearer " prefix

	newToken, err := h.jwtService.RefreshToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to refresh token",
		})
	}

	return c.JSON(fiber.Map{
		"token": newToken,
	})
}

func generateRandomState() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
