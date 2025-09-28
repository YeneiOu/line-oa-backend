package entities

import (
	"errors"
	"time"
)

// BookingStatus represents the status of a booking
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

// Booking represents the core booking entity
type Booking struct {
	ID          string
	UserID      string
	ServiceName string
	BookingDate time.Time
	Notes       string
	Status      BookingStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewBooking creates a new booking entity
func NewBooking(userID, serviceName string, bookingDate time.Time, notes string) *Booking {
	now := time.Now()
	return &Booking{
		UserID:      userID,
		ServiceName: serviceName,
		BookingDate: bookingDate,
		Notes:       notes,
		Status:      BookingStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Confirm confirms the booking
func (b *Booking) Confirm() error {
	if b.Status != BookingStatusPending {
		return errors.New("only pending bookings can be confirmed")
	}
	b.Status = BookingStatusConfirmed
	b.UpdatedAt = time.Now()
	return nil
}

// Cancel cancels the booking
func (b *Booking) Cancel() error {
	if b.Status == BookingStatusCompleted {
		return errors.New("completed bookings cannot be cancelled")
	}
	b.Status = BookingStatusCancelled
	b.UpdatedAt = time.Now()
	return nil
}

// Complete marks the booking as completed
func (b *Booking) Complete() error {
	if b.Status != BookingStatusConfirmed {
		return errors.New("only confirmed bookings can be completed")
	}
	b.Status = BookingStatusCompleted
	b.UpdatedAt = time.Now()
	return nil
}

// UpdateNotes updates the booking notes
func (b *Booking) UpdateNotes(notes string) {
	b.Notes = notes
	b.UpdatedAt = time.Now()
}

// IsValid validates the booking entity
func (b *Booking) IsValid() bool {
	return b.UserID != "" && b.ServiceName != "" && !b.BookingDate.IsZero()
}

// CanBeModified checks if the booking can be modified
func (b *Booking) CanBeModified() bool {
	return b.Status == BookingStatusPending || b.Status == BookingStatusConfirmed
}
