package handlers

import (
	"context"
	"line-oa-backend/internal/application"
	"time"

	"github.com/gofiber/fiber/v2"
)

// BookingHandlerAdapter implements BookingHandler using Fiber
type BookingHandlerAdapter struct {
	bookingService *application.BookingService
}

// NewBookingHandlerAdapter creates a new booking handler adapter
func NewBookingHandlerAdapter(bookingService *application.BookingService) *BookingHandlerAdapter {
	return &BookingHandlerAdapter{
		bookingService: bookingService,
	}
}

type CreateBookingRequest struct {
	ServiceName string    `json:"service_name" validate:"required"`
	BookingDate time.Time `json:"booking_date" validate:"required"`
	Notes       string    `json:"notes"`
}

type UpdateBookingRequest struct {
	ServiceName string    `json:"service_name" validate:"required"`
	BookingDate time.Time `json:"booking_date" validate:"required"`
	Notes       string    `json:"notes"`
}

// CreateBooking creates a new booking
func (h *BookingHandlerAdapter) CreateBooking(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var req CreateBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.ServiceName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Service name is required",
		})
	}

	if req.BookingDate.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Booking date is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	booking, err := h.bookingService.CreateBooking(ctx, userID, req.ServiceName, req.BookingDate, req.Notes)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(booking)
}

// GetBookings retrieves all bookings for the authenticated user
func (h *BookingHandlerAdapter) GetBookings(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bookings, err := h.bookingService.GetUserBookings(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve bookings",
		})
	}

	return c.JSON(bookings)
}

// GetBooking retrieves a specific booking
func (h *BookingHandlerAdapter) GetBooking(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	bookingID := c.Params("id")
	if bookingID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Booking ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	booking, err := h.bookingService.GetBooking(ctx, bookingID, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Booking not found",
		})
	}

	return c.JSON(booking)
}

// UpdateBooking updates a booking
func (h *BookingHandlerAdapter) UpdateBooking(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	bookingID := c.Params("id")
	if bookingID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Booking ID is required",
		})
	}

	var req UpdateBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.ServiceName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Service name is required",
		})
	}

	if req.BookingDate.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Booking date is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	booking, err := h.bookingService.UpdateBooking(ctx, bookingID, userID, req.ServiceName, req.BookingDate, req.Notes)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(booking)
}

// CancelBooking cancels a booking
func (h *BookingHandlerAdapter) CancelBooking(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	bookingID := c.Params("id")
	if bookingID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Booking ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.bookingService.CancelBooking(ctx, bookingID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Booking cancelled successfully",
	})
}
