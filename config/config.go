package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// MongoDB Database
	MongoURI      string
	MongoDatabase string

	// JWT
	JWTSecret string

	// LINE OAuth
	LINEChannelID     string
	LINEChannelSecret string
	LINERedirectURI   string

	// LINE Messaging API
	LINEChannelAccessToken string

	// Server
	Port        string
	FrontendURL string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		// MongoDB Database
		MongoURI:      os.Getenv("MONGO_URI"),
		MongoDatabase: os.Getenv("MONGO_DATABASE"),

		// JWT
		JWTSecret: os.Getenv("JWT_SECRET"),

		// LINE OAuth
		LINEChannelID:     os.Getenv("LINE_CHANNEL_ID"),
		LINEChannelSecret: os.Getenv("LINE_CHANNEL_SECRET"),
		LINERedirectURI:   os.Getenv("LINE_REDIRECT_URI"),

		// LINE Messaging API
		LINEChannelAccessToken: os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),

		// Server
		Port:        os.Getenv("PORT"),
		FrontendURL: os.Getenv("FRONTEND_URL"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
