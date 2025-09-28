package entities

import (
	"time"
)

// User represents the core user entity
type User struct {
	ID         string
	LineUserID string
	Name       string
	Email      string
	PictureURL string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewUser creates a new user entity
func NewUser(lineUserID, name, email, pictureURL string) *User {
	now := time.Now()
	return &User{
		LineUserID: lineUserID,
		Name:       name,
		Email:      email,
		PictureURL: pictureURL,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(name, email, pictureURL string) {
	u.Name = name
	u.Email = email
	u.PictureURL = pictureURL
	u.UpdatedAt = time.Now()
}

// IsValid validates the user entity
func (u *User) IsValid() bool {
	return u.LineUserID != "" && u.Name != ""
}
