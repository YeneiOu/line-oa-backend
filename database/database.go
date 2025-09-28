package database

import (
	"context"
	"log"
	"time"

	"line-oa-backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func Connect(cfg *config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Test the connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	Client = client
	Database = client.Database(cfg.MongoDatabase)

	log.Println("MongoDB connected successfully")
	log.Printf("Using database: %s", cfg.MongoDatabase)

	// Create indexes for better performance
	createIndexes()
}

func GetDatabase() *mongo.Database {
	return Database
}

func GetClient() *mongo.Client {
	return Client
}

func createIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create unique index on line_user_id for users collection
	usersCollection := Database.Collection("users")
	indexModel := mongo.IndexModel{
		Keys:    map[string]int{"line_user_id": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := usersCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Warning: Failed to create index on users collection: %v", err)
	}

	// Create index on user_id for bookings collection
	bookingsCollection := Database.Collection("bookings")
	bookingIndexModel := mongo.IndexModel{
		Keys: map[string]int{"user_id": 1},
	}

	_, err = bookingsCollection.Indexes().CreateOne(ctx, bookingIndexModel)
	if err != nil {
		log.Printf("Warning: Failed to create index on bookings collection: %v", err)
	}

	log.Println("Database indexes created successfully")
}

func Disconnect() {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		if err := Client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		} else {
			log.Println("MongoDB connection closed")
		}
	}
}
