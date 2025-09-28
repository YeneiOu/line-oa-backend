package valueobjects

import "time"

// AuthToken represents authentication token information
type AuthToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	TokenType    string
	Scope        string
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID     string `json:"user_id"`
	LineUserID string `json:"line_user_id"`
	ExpiresAt  int64  `json:"exp"`
	IssuedAt   int64  `json:"iat"`
}

// NewJWTClaims creates new JWT claims
func NewJWTClaims(userID, lineUserID string, expirationTime time.Time) *JWTClaims {
	now := time.Now()
	return &JWTClaims{
		UserID:     userID,
		LineUserID: lineUserID,
		ExpiresAt:  expirationTime.Unix(),
		IssuedAt:   now.Unix(),
	}
}

// IsExpired checks if the token is expired
func (c *JWTClaims) IsExpired() bool {
	return time.Now().Unix() > c.ExpiresAt
}
