package services

import (
	"errors"
	"line-oa-backend/internal/domain/valueobjects"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTServiceAdapter implements JWTService
type JWTServiceAdapter struct {
	secretKey []byte
}

// NewJWTServiceAdapter creates a new JWT service adapter
func NewJWTServiceAdapter(secretKey string) *JWTServiceAdapter {
	return &JWTServiceAdapter{
		secretKey: []byte(secretKey),
	}
}

// GenerateToken generates a JWT token
func (s *JWTServiceAdapter) GenerateToken(userID, lineUserID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := valueobjects.NewJWTClaims(userID, lineUserID, expirationTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":      claims.UserID,
		"line_user_id": claims.LineUserID,
		"exp":          claims.ExpiresAt,
		"iat":          claims.IssuedAt,
	})

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns claims
func (s *JWTServiceAdapter) ValidateToken(tokenString string) (*valueobjects.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}

	lineUserID, ok := claims["line_user_id"].(string)
	if !ok {
		return nil, errors.New("invalid line_user_id in token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("invalid exp in token")
	}

	iat, ok := claims["iat"].(float64)
	if !ok {
		return nil, errors.New("invalid iat in token")
	}

	jwtClaims := &valueobjects.JWTClaims{
		UserID:     userID,
		LineUserID: lineUserID,
		ExpiresAt:  int64(exp),
		IssuedAt:   int64(iat),
	}

	if jwtClaims.IsExpired() {
		return nil, errors.New("token expired")
	}

	return jwtClaims, nil
}

// RefreshToken refreshes a JWT token
func (s *JWTServiceAdapter) RefreshToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		// Allow expired tokens for refresh
		if err.Error() != "token expired" {
			return "", err
		}
		
		// Parse token without validation to get claims
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return s.secretKey, nil
		})
		if err != nil {
			return "", err
		}

		mapClaims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return "", errors.New("invalid token claims")
		}

		userID, ok := mapClaims["user_id"].(string)
		if !ok {
			return "", errors.New("invalid user_id in token")
		}

		lineUserID, ok := mapClaims["line_user_id"].(string)
		if !ok {
			return "", errors.New("invalid line_user_id in token")
		}

		return s.GenerateToken(userID, lineUserID)
	}

	// Generate new token with same claims
	return s.GenerateToken(claims.UserID, claims.LineUserID)
}
