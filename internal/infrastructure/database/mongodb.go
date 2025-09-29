package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

// Connect establishes connection to MongoDB
func Connect(mongoURI, databaseName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Configure client options for MongoDB Atlas
	clientOptions := options.Client().ApplyURI(mongoURI)
	
	// Add additional options for Atlas connections
	clientOptions.SetMaxPoolSize(10)
	clientOptions.SetMinPoolSize(1)
	clientOptions.SetMaxConnIdleTime(30 * time.Second)
	
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Failed to create MongoDB client: %v", err)
		return nil, err
	}

	// Test the connection with a longer timeout for Atlas
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer pingCancel()
	
	err = client.Ping(pingCtx, nil)
	if err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		return nil, err
	}

	db = client.Database(databaseName)
	log.Printf("Successfully connected to MongoDB database: %s", databaseName)
	
	return db, nil
}

// GetDatabase returns the database instance
func GetDatabase() *mongo.Database {
	return db
}
