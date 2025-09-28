package ports

import (
	"context"
	"line-oa-backend/internal/domain/entities"
)

// UserRepository defines the interface for user data persistence
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByLineUserID(ctx context.Context, lineUserID string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id string) error
}

// BookingRepository defines the interface for booking data persistence
type BookingRepository interface {
	Create(ctx context.Context, booking *entities.Booking) error
	GetByID(ctx context.Context, id string) (*entities.Booking, error)
	GetByUserID(ctx context.Context, userID string) ([]*entities.Booking, error)
	Update(ctx context.Context, booking *entities.Booking) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]*entities.Booking, error)
}
