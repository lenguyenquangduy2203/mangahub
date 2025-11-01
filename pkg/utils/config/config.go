package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	API_PORT   string
	DB_PATH    string
	JWT_SECRET string
}

// LoadConfig initializes the Config struct.
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — using system environment variables.")
	}

	cfg := &Config{
		API_PORT:   getEnv("API_PORT", "3000"),              // Fallback to 3000
		DB_PATH:    getEnv("DB_PATH", "./data/mangahub.db"), // Fallback path
		JWT_SECRET: getEnv("JWT_SECRET", "super-secret-key"),
	}

	// Example: check for a critical setting
	if cfg.JWT_SECRET == "super-secret-key" {
		log.Println("WARNING: Using default JWT_SECRET. Update .env for production.")
	}

	return cfg, nil
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
