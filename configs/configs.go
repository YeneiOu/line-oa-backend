package configs

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
	Port string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		// MongoDB Database
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGO_DATABASE", "line_oa_backend"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "default-secret-key"),

		// LINE OAuth
		LINEChannelID:     getEnv("LINE_CHANNEL_ID", ""),
		LINEChannelSecret: getEnv("LINE_CHANNEL_SECRET", ""),
		LINERedirectURI:   getEnv("LINE_REDIRECT_URI", "http://localhost:3000/callback"),

		// LINE Messaging API
		LINEChannelAccessToken: getEnv("LINE_CHANNEL_ACCESS_TOKEN", ""),

		// Server
		Port: getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
