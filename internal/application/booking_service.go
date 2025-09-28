package application

import (
	"context"
	"errors"
	"line-oa-backend/internal/domain/entities"
	"line-oa-backend/internal/ports"
	"time"
)

// BookingService handles booking use cases
type BookingService struct {
	bookingRepo         ports.BookingRepository
	userRepo            ports.UserRepository
	notificationService ports.NotificationService
}

// NewBookingService creates a new booking service
func NewBookingService(
	bookingRepo ports.BookingRepository,
	userRepo ports.UserRepository,
	notificationService ports.NotificationService,
) *BookingService {
	return &BookingService{
		bookingRepo:         bookingRepo,
		userRepo:            userRepo,
		notificationService: notificationService,
	}
}

// CreateBooking creates a new booking
func (s *BookingService) CreateBooking(ctx context.Context, userID, serviceName string, bookingDate time.Time, notes string) (*entities.Booking, error) {
	// Validate user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Validate booking date is in the future
	if bookingDate.Before(time.Now()) {
		return nil, errors.New("booking date must be in the future")
	}

	// Create booking entity
	booking := entities.NewBooking(userID, serviceName, bookingDate, notes)
	if !booking.IsValid() {
		return nil, errors.New("invalid booking data")
	}

	// Save booking
	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		return nil, err
	}

	// Send notification (async)
	go func() {
		if err := s.notificationService.SendBookingConfirmation(userID, booking.ID); err != nil {
			// Log error but don't fail the booking creation
		}
	}()

	return booking, nil
}

// GetBooking retrieves a booking by ID
func (s *BookingService) GetBooking(ctx context.Context, bookingID, userID string) (*entities.Booking, error) {
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	// Ensure user owns the booking
	if booking.UserID != userID {
		return nil, errors.New("booking not found")
	}

	return booking, nil
}

// GetUserBookings retrieves all bookings for a user
func (s *BookingService) GetUserBookings(ctx context.Context, userID string) ([]*entities.Booking, error) {
	return s.bookingRepo.GetByUserID(ctx, userID)
}

// UpdateBooking updates a booking
func (s *BookingService) UpdateBooking(ctx context.Context, bookingID, userID string, serviceName string, bookingDate time.Time, notes string) (*entities.Booking, error) {
	booking, err := s.GetBooking(ctx, bookingID, userID)
	if err != nil {
		return nil, err
	}

	if !booking.CanBeModified() {
		return nil, errors.New("booking cannot be modified")
	}

	// Validate new booking date
	if bookingDate.Before(time.Now()) {
		return nil, errors.New("booking date must be in the future")
	}

	// Update booking
	booking.ServiceName = serviceName
	booking.BookingDate = bookingDate
	booking.UpdateNotes(notes)

	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, err
	}

	return booking, nil
}

// ConfirmBooking confirms a booking
func (s *BookingService) ConfirmBooking(ctx context.Context, bookingID, userID string) (*entities.Booking, error) {
	booking, err := s.GetBooking(ctx, bookingID, userID)
	if err != nil {
		return nil, err
	}

	if err := booking.Confirm(); err != nil {
		return nil, err
	}

	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, err
	}

	return booking, nil
}

// CancelBooking cancels a booking
func (s *BookingService) CancelBooking(ctx context.Context, bookingID, userID string) error {
	booking, err := s.GetBooking(ctx, bookingID, userID)
	if err != nil {
		return err
	}

	if err := booking.Cancel(); err != nil {
		return err
	}

	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return err
	}

	// Send cancellation notification
	go func() {
		if err := s.notificationService.SendBookingCancellation(userID, booking.ID); err != nil {
			// Log error but don't fail the cancellation
		}
	}()

	return nil
}

// CompleteBooking marks a booking as completed
func (s *BookingService) CompleteBooking(ctx context.Context, bookingID, userID string) (*entities.Booking, error) {
	booking, err := s.GetBooking(ctx, bookingID, userID)
	if err != nil {
		return nil, err
	}

	if err := booking.Complete(); err != nil {
		return nil, err
	}

	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, err
	}

	return booking, nil
}
