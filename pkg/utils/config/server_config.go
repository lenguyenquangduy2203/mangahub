package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	API_CONFIG    APIConfig
	SOCKET_CONFIG SocketConfig
}

type APIConfig struct {
	API_PORT                  string
	DB_PATH                   string
	JWT_SECRET                string
	JWT_ACCESS_TOKEN_LIFETIME string
}

type SocketConfig struct {
	CHAT_ROOMS []string
}

// LoadConfig initializes the Config struct.
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — using system environment variables.")
	}

	cfg := &Config{
		API_CONFIG: APIConfig{
			API_PORT:                  getEnvAsStr("API_PORT", "3000"),              // Fallback to 3000
			DB_PATH:                   getEnvAsStr("DB_PATH", "./data/mangahub.db"), // Fallback path
			JWT_SECRET:                getEnvAsStr("JWT_SECRET", "super-secret-key"),
			JWT_ACCESS_TOKEN_LIFETIME: getEnvAsStr("JWT_ACCESS_TOKEN_LIFETIME", "4h"),
		},

		SOCKET_CONFIG: SocketConfig{
			CHAT_ROOMS: getEnvAsListStr("CHAT_ROOMS", "general"),
		},
	}

	// Example: check for a critical setting
	if cfg.API_CONFIG.JWT_SECRET == "super-secret-key" {
		log.Println("WARNING: Using default JWT_SECRET. Update .env for production.")
	}

	return cfg, nil
}

func getEnvAsStr(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsListStr(key, fallback string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, ",")
	}

	return []string{fallback}
}
