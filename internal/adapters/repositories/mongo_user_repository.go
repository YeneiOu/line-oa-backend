package repositories

import (
	"context"
	"line-oa-backend/internal/domain/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoUserRepository implements UserRepository using MongoDB
type MongoUserRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRepository creates a new MongoDB user repository
func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		collection: db.Collection("users"),
	}
}

// mongoUser represents the MongoDB document structure
type mongoUser struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	LineUserID string             `bson:"line_user_id"`
	Name       string             `bson:"name"`
	Email      string             `bson:"email"`
	PictureURL string             `bson:"picture_url"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

// toEntity converts MongoDB document to domain entity
func (mu *mongoUser) toEntity() *entities.User {
	return &entities.User{
		ID:         mu.ID.Hex(),
		LineUserID: mu.LineUserID,
		Name:       mu.Name,
		Email:      mu.Email,
		PictureURL: mu.PictureURL,
		CreatedAt:  mu.CreatedAt,
		UpdatedAt:  mu.UpdatedAt,
	}
}

// fromEntity converts domain entity to MongoDB document
func fromUserEntity(user *entities.User) *mongoUser {
	mu := &mongoUser{
		LineUserID: user.LineUserID,
		Name:       user.Name,
		Email:      user.Email,
		PictureURL: user.PictureURL,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	if user.ID != "" {
		if objID, err := primitive.ObjectIDFromHex(user.ID); err == nil {
			mu.ID = objID
		}
	}

	return mu
}

// Create creates a new user
func (r *MongoUserRepository) Create(ctx context.Context, user *entities.User) error {
	mongoUser := fromUserEntity(user)
	mongoUser.ID = primitive.NewObjectID()
	
	result, err := r.collection.InsertOne(ctx, mongoUser)
	if err != nil {
		return err
	}
	
	user.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

// GetByID retrieves a user by ID
func (r *MongoUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var mongoUser mongoUser
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&mongoUser)
	if err != nil {
		return nil, err
	}

	return mongoUser.toEntity(), nil
}

// GetByLineUserID retrieves a user by LINE user ID
func (r *MongoUserRepository) GetByLineUserID(ctx context.Context, lineUserID string) (*entities.User, error) {
	var mongoUser mongoUser
	err := r.collection.FindOne(ctx, bson.M{"line_user_id": lineUserID}).Decode(&mongoUser)
	if err != nil {
		return nil, err
	}

	return mongoUser.toEntity(), nil
}

// Update updates a user
func (r *MongoUserRepository) Update(ctx context.Context, user *entities.User) error {
	objID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name":        user.Name,
			"email":       user.Email,
			"picture_url": user.PictureURL,
			"updated_at":  user.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

// Delete deletes a user
func (r *MongoUserRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
