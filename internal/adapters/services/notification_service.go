package services

import (
	"line-oa-backend/internal/ports"
)

// NotificationServiceAdapter implements NotificationService
type NotificationServiceAdapter struct {
	lineMessaging ports.LINEMessagingService
}

// NewNotificationServiceAdapter creates a new notification service adapter
func NewNotificationServiceAdapter(lineMessaging ports.LINEMessagingService) *NotificationServiceAdapter {
	return &NotificationServiceAdapter{
		lineMessaging: lineMessaging,
	}
}

// SendBookingConfirmation sends booking confirmation notification
func (s *NotificationServiceAdapter) SendBookingConfirmation(userID string, bookingID string) error {
	message := "🎉 Your booking has been confirmed! We'll send you a reminder closer to your appointment date."
	return s.lineMessaging.SendTextMessage(userID, message)
}

// SendBookingReminder sends booking reminder notification
func (s *NotificationServiceAdapter) SendBookingReminder(userID string, bookingID string) error {
	message := "⏰ Reminder: You have an upcoming appointment tomorrow. Please arrive 10 minutes early."
	return s.lineMessaging.SendTextMessage(userID, message)
}

// SendBookingCancellation sends booking cancellation notification
func (s *NotificationServiceAdapter) SendBookingCancellation(userID string, bookingID string) error {
	message := "❌ Your booking has been cancelled. If you need to reschedule, please make a new booking."
	return s.lineMessaging.SendTextMessage(userID, message)
}
