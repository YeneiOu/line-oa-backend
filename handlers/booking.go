package handlers

import (
	"context"
	"log"
	"time"

	"line-oa-backend/database"
	"line-oa-backend/models"
	"line-oa-backend/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookingHandler struct {
	lineMessaging *services.LINEMessagingService
}

type CreateBookingResponse struct {
	Booking models.Booking `json:"booking"`
	Message string         `json:"message"`
}

func NewBookingHandler(lineMessaging *services.LINEMessagingService) *BookingHandler {
	return &BookingHandler{
		lineMessaging: lineMessaging,
	}
}

// CreateBooking creates a new booking and sends LINE notification
func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID)
	lineUserID := c.Locals("lineUserID").(string)

	var req models.BookingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
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

	// Check if booking date is in the future
	if req.BookingDate.Before(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Booking date must be in the future",
		})
	}

	// Create booking
	booking := models.Booking{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		ServiceName: req.ServiceName,
		BookingDate: req.BookingDate,
		Notes:       req.Notes,
		Status:      "confirmed",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	db := database.GetDatabase()
	bookingsCollection := db.Collection("bookings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := bookingsCollection.InsertOne(ctx, booking)
	if err != nil {
		log.Printf("Failed to create booking: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create booking",
		})
	}

	// Send LINE notification
	if err := h.lineMessaging.SendBookingConfirmation(lineUserID, &booking); err != nil {
		log.Printf("Failed to send LINE notification: %v", err)
		// Don't return error here as booking was created successfully
		// Just log the error and continue
	}

	return c.Status(fiber.StatusCreated).JSON(CreateBookingResponse{
		Booking: booking,
		Message: "Booking created successfully and notification sent",
	})
}

// GetBookings returns user's bookings
func (h *BookingHandler) GetBookings(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID)

	db := database.GetDatabase()
	bookingsCollection := db.Collection("bookings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find bookings for the user, sorted by booking_date descending
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "booking_date", Value: -1}})

	cursor, err := bookingsCollection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("Failed to get bookings: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve bookings",
		})
	}
	defer cursor.Close(ctx)

	var bookings []models.Booking
	if err = cursor.All(ctx, &bookings); err != nil {
		log.Printf("Failed to decode bookings: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve bookings",
		})
	}

	return c.JSON(fiber.Map{
		"bookings": bookings,
	})
}

// GetBooking returns a specific booking
func (h *BookingHandler) GetBooking(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID)
	bookingIDStr := c.Params("id")

	bookingID, err := primitive.ObjectIDFromHex(bookingIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking ID",
		})
	}

	db := database.GetDatabase()
	bookingsCollection := db.Collection("bookings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var booking models.Booking
	filter := bson.M{"_id": bookingID, "user_id": userID}
	err = bookingsCollection.FindOne(ctx, filter).Decode(&booking)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	return c.JSON(booking)
}

// UpdateBooking updates a booking
func (h *BookingHandler) UpdateBooking(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID)
	lineUserID := c.Locals("lineUserID").(string)
	bookingIDStr := c.Params("id")

	bookingID, err := primitive.ObjectIDFromHex(bookingIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking ID",
		})
	}

	var req models.BookingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	db := database.GetDatabase()
	bookingsCollection := db.Collection("bookings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var booking models.Booking

	// Find existing booking
	filter := bson.M{"_id": bookingID, "user_id": userID}
	err = bookingsCollection.FindOne(ctx, filter).Decode(&booking)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Prepare update document
	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}

	// Update booking fields
	if req.ServiceName != "" {
		update["$set"].(bson.M)["service_name"] = req.ServiceName
		booking.ServiceName = req.ServiceName
	}
	if !req.BookingDate.IsZero() {
		if req.BookingDate.Before(time.Now()) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Booking date must be in the future",
			})
		}
		update["$set"].(bson.M)["booking_date"] = req.BookingDate
		booking.BookingDate = req.BookingDate
	}
	update["$set"].(bson.M)["notes"] = req.Notes
	booking.Notes = req.Notes
	booking.UpdatedAt = time.Now()

	_, err = bookingsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Failed to update booking: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update booking",
		})
	}

	// Send update notification
	updateMessage := "📝 การจองของคุณได้รับการอัปเดตแล้ว!\n\n" +
		"บริการ: " + booking.ServiceName + "\n" +
		"วันที่: " + booking.BookingDate.Format("2 January 2006 15:04")

	if err := h.lineMessaging.SendCustomMessage(lineUserID, updateMessage); err != nil {
		log.Printf("Failed to send update notification: %v", err)
	}

	return c.JSON(fiber.Map{
		"booking": booking,
		"message": "Booking updated successfully",
	})
}

// CancelBooking cancels a booking
func (h *BookingHandler) CancelBooking(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID)
	lineUserID := c.Locals("lineUserID").(string)
	bookingIDStr := c.Params("id")

	bookingID, err := primitive.ObjectIDFromHex(bookingIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking ID",
		})
	}

	db := database.GetDatabase()
	bookingsCollection := db.Collection("bookings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var booking models.Booking

	// Find existing booking
	filter := bson.M{"_id": bookingID, "user_id": userID}
	err = bookingsCollection.FindOne(ctx, filter).Decode(&booking)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Update status to cancelled
	update := bson.M{
		"$set": bson.M{
			"status":     "cancelled",
			"updated_at": time.Now(),
		},
	}

	_, err = bookingsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Failed to cancel booking: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to cancel booking",
		})
	}

	// Update local booking object
	booking.Status = "cancelled"
	booking.UpdatedAt = time.Now()

	// Send cancellation notification
	cancelMessage := "❌ การจองของคุณได้รับการยกเลิกแล้ว\n\n" +
		"บริการ: " + booking.ServiceName + "\n" +
		"วันที่: " + booking.BookingDate.Format("2 January 2006 15:04")

	if err := h.lineMessaging.SendCustomMessage(lineUserID, cancelMessage); err != nil {
		log.Printf("Failed to send cancellation notification: %v", err)
	}

	return c.JSON(fiber.Map{
		"booking": booking,
		"message": "Booking cancelled successfully",
	})
}
