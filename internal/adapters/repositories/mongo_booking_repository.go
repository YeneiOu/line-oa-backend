package repositories

import (
	"context"
	"line-oa-backend/internal/domain/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoBookingRepository implements BookingRepository using MongoDB
type MongoBookingRepository struct {
	collection *mongo.Collection
}

// NewMongoBookingRepository creates a new MongoDB booking repository
func NewMongoBookingRepository(db *mongo.Database) *MongoBookingRepository {
	return &MongoBookingRepository{
		collection: db.Collection("bookings"),
	}
}

// mongoBooking represents the MongoDB document structure
type mongoBooking struct {
	ID          primitive.ObjectID    `bson:"_id,omitempty"`
	UserID      primitive.ObjectID    `bson:"user_id"`
	ServiceName string                `bson:"service_name"`
	BookingDate time.Time             `bson:"booking_date"`
	Notes       string                `bson:"notes"`
	Status      entities.BookingStatus `bson:"status"`
	CreatedAt   time.Time             `bson:"created_at"`
	UpdatedAt   time.Time             `bson:"updated_at"`
}

// toEntity converts MongoDB document to domain entity
func (mb *mongoBooking) toEntity() *entities.Booking {
	return &entities.Booking{
		ID:          mb.ID.Hex(),
		UserID:      mb.UserID.Hex(),
		ServiceName: mb.ServiceName,
		BookingDate: mb.BookingDate,
		Notes:       mb.Notes,
		Status:      mb.Status,
		CreatedAt:   mb.CreatedAt,
		UpdatedAt:   mb.UpdatedAt,
	}
}

// fromEntity converts domain entity to MongoDB document
func fromBookingEntity(booking *entities.Booking) (*mongoBooking, error) {
	userObjID, err := primitive.ObjectIDFromHex(booking.UserID)
	if err != nil {
		return nil, err
	}

	mb := &mongoBooking{
		UserID:      userObjID,
		ServiceName: booking.ServiceName,
		BookingDate: booking.BookingDate,
		Notes:       booking.Notes,
		Status:      booking.Status,
		CreatedAt:   booking.CreatedAt,
		UpdatedAt:   booking.UpdatedAt,
	}

	if booking.ID != "" {
		if objID, err := primitive.ObjectIDFromHex(booking.ID); err == nil {
			mb.ID = objID
		}
	}

	return mb, nil
}

// Create creates a new booking
func (r *MongoBookingRepository) Create(ctx context.Context, booking *entities.Booking) error {
	mongoBooking, err := fromBookingEntity(booking)
	if err != nil {
		return err
	}
	
	mongoBooking.ID = primitive.NewObjectID()
	
	result, err := r.collection.InsertOne(ctx, mongoBooking)
	if err != nil {
		return err
	}
	
	booking.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

// GetByID retrieves a booking by ID
func (r *MongoBookingRepository) GetByID(ctx context.Context, id string) (*entities.Booking, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var mongoBooking mongoBooking
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&mongoBooking)
	if err != nil {
		return nil, err
	}

	return mongoBooking.toEntity(), nil
}

// GetByUserID retrieves all bookings for a user
func (r *MongoBookingRepository) GetByUserID(ctx context.Context, userID string) ([]*entities.Booking, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userObjID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []*entities.Booking
	for cursor.Next(ctx) {
		var mongoBooking mongoBooking
		if err := cursor.Decode(&mongoBooking); err != nil {
			return nil, err
		}
		bookings = append(bookings, mongoBooking.toEntity())
	}

	return bookings, cursor.Err()
}

// Update updates a booking
func (r *MongoBookingRepository) Update(ctx context.Context, booking *entities.Booking) error {
	objID, err := primitive.ObjectIDFromHex(booking.ID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"service_name":  booking.ServiceName,
			"booking_date":  booking.BookingDate,
			"notes":         booking.Notes,
			"status":        booking.Status,
			"updated_at":    booking.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

// Delete deletes a booking
func (r *MongoBookingRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// List retrieves bookings with pagination
func (r *MongoBookingRepository) List(ctx context.Context, offset, limit int) ([]*entities.Booking, error) {
	opts := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []*entities.Booking
	for cursor.Next(ctx) {
		var mongoBooking mongoBooking
		if err := cursor.Decode(&mongoBooking); err != nil {
			return nil, err
		}
		bookings = append(bookings, mongoBooking.toEntity())
	}

	return bookings, cursor.Err()
}
