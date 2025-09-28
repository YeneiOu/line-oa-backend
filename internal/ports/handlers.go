package ports

import "github.com/gofiber/fiber/v2"

// AuthHandler defines the interface for authentication HTTP handlers
type AuthHandler interface {
	Login(c *fiber.Ctx) error
	Callback(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
}

// BookingHandler defines the interface for booking HTTP handlers
type BookingHandler interface {
	CreateBooking(c *fiber.Ctx) error
	GetBookings(c *fiber.Ctx) error
	GetBooking(c *fiber.Ctx) error
	UpdateBooking(c *fiber.Ctx) error
	CancelBooking(c *fiber.Ctx) error
}
