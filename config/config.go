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
	// Try to load .env file from current directory and parent directories
	envFiles := []string{".env", "../.env", "../../.env"}
	loaded := false
	
	for _, envFile := range envFiles {
		if err := godotenv.Load(envFile); err == nil {
			log.Printf("Loaded environment from: %s", envFile)
			loaded = true
			break
		}
	}
	
	if !loaded {
		log.Println("No .env file found, using system environment variables")
	}

	config := &Config{
		// MongoDB Database
		MongoURI:      getEnv("MONGO_URI", ""),
		MongoDatabase: getEnv("MONGO_DATABASE", "line_oa_backend"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "default-secret-key-change-in-production"),

		// LINE OAuth
		LINEChannelID:     getEnv("LINE_CHANNEL_ID", ""),
		LINEChannelSecret: getEnv("LINE_CHANNEL_SECRET", ""),
		LINERedirectURI:   getEnv("LINE_REDIRECT_URI", "http://localhost:3000/callback"),

		// LINE Messaging API
		LINEChannelAccessToken: getEnv("LINE_CHANNEL_ACCESS_TOKEN", ""),

		// Server
		Port:        getEnv("PORT", "8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
	}

	// Validate required fields
	if config.MongoURI == "" {
		log.Fatal("MONGO_URI environment variable is required")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
