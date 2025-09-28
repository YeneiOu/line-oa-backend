package application

import (
	"context"
	"errors"
	"line-oa-backend/internal/domain/entities"
	"line-oa-backend/internal/ports"
)

// AuthService handles authentication use cases
type AuthService struct {
	userRepo         ports.UserRepository
	lineOAuthService ports.LINEOAuthService
	jwtService       ports.JWTService
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo ports.UserRepository,
	lineOAuthService ports.LINEOAuthService,
	jwtService ports.JWTService,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		lineOAuthService: lineOAuthService,
		jwtService:       jwtService,
	}
}

// GetAuthURL generates LINE OAuth URL
func (s *AuthService) GetAuthURL(state string) string {
	return s.lineOAuthService.GetAuthURL(state)
}

// AuthenticateWithLINE handles LINE OAuth callback and user authentication
func (s *AuthService) AuthenticateWithLINE(ctx context.Context, code string) (*entities.User, string, error) {
	// Exchange code for token
	tokenResp, err := s.lineOAuthService.ExchangeCodeForToken(code)
	if err != nil {
		return nil, "", err
	}

	// Get user profile from LINE
	profile, err := s.lineOAuthService.GetProfile(tokenResp.AccessToken)
	if err != nil {
		return nil, "", err
	}

	if !profile.IsValid() {
		return nil, "", errors.New("invalid LINE profile")
	}

	// Find or create user
	user, err := s.userRepo.GetByLineUserID(ctx, profile.UserID)
	if err != nil {
		// User doesn't exist, create new one
		user = entities.NewUser(profile.UserID, profile.DisplayName, "", profile.PictureURL)
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, "", err
		}
	} else {
		// Update existing user profile
		user.UpdateProfile(profile.DisplayName, user.Email, profile.PictureURL)
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, "", err
		}
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.LineUserID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GetUser retrieves user by ID
func (s *AuthService) GetUser(ctx context.Context, userID string) (*entities.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// RefreshToken refreshes JWT token
func (s *AuthService) RefreshToken(tokenString string) (string, error) {
	return s.jwtService.RefreshToken(tokenString)
}

// ValidateToken validates JWT token
func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
