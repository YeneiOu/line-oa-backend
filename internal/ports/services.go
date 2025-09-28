package ports

import (
	"line-oa-backend/internal/domain/valueobjects"
)

// LINEOAuthService defines the interface for LINE OAuth operations
type LINEOAuthService interface {
	GetAuthURL(state string) string
	ExchangeCodeForToken(code string) (*valueobjects.AuthToken, error)
	GetProfile(accessToken string) (*valueobjects.LINEProfile, error)
}

// LINEMessagingService defines the interface for LINE Messaging API operations
type LINEMessagingService interface {
	SendTextMessage(userID, message string) error
	SendFlexMessage(userID string, flexMessage interface{}) error
	BroadcastMessage(message string) error
}

// JWTService defines the interface for JWT token operations
type JWTService interface {
	GenerateToken(userID, lineUserID string) (string, error)
	ValidateToken(tokenString string) (*valueobjects.JWTClaims, error)
	RefreshToken(tokenString string) (string, error)
}

// NotificationService defines the interface for sending notifications
type NotificationService interface {
	SendBookingConfirmation(userID string, bookingID string) error
	SendBookingReminder(userID string, bookingID string) error
	SendBookingCancellation(userID string, bookingID string) error
}
