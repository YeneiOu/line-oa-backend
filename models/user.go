package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	LineUserID string             `json:"line_user_id" bson:"line_user_id"`
	Name       string             `json:"name" bson:"name"`
	Email      string             `json:"email" bson:"email"`
	PictureURL string             `json:"picture_url" bson:"picture_url"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

type Booking struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	ServiceName string             `json:"service_name" bson:"service_name"`
	BookingDate time.Time          `json:"booking_date" bson:"booking_date"`
	Notes       string             `json:"notes" bson:"notes"`
	Status      string             `json:"status" bson:"status"` // confirmed, cancelled, completed
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// BookingRequest represents the request payload for creating a booking
type BookingRequest struct {
	ServiceName string    `json:"service_name" validate:"required"`
	BookingDate time.Time `json:"booking_date" validate:"required"`
	Notes       string    `json:"notes"`
}
