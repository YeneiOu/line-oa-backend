package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the core user entity
type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	LineUserID string             `bson:"line_user_id" json:"line_user_id"`
	Name       string             `bson:"name" json:"name"`
	Email      string             `bson:"email" json:"email"`
	PictureURL string             `bson:"picture_url" json:"picture_url"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
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

// Responses represents the standard API response structure
type Responses struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
