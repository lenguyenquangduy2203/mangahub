package config

import (
	"log"
	"os"
	"strconv"

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
	// MaxConnections    int
	// ReadBufferSize    int
	// WriteBufferSize   int
	// HandshakeTimeout  time.Duration
	// PongWait          time.Duration
	// PingPeriod        time.Duration
	// WriteWait         time.Duration
	// MaxMessageSize    int64
	// EnableCompression bool
	MAX_CONNECTIONS   int
	READ_BUFFER_SIZE  int // As byte
	WRITE_BUFFER_SIZE int // As byte
	PONG_WAIT         int // As second
	PING_PERIOD       int // AS second
	WRITE_WAIT        int // As second
	MAX_MESSAGE       int
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
			MAX_CONNECTIONS:   getEnvAsInt("MAX_CONNECTIONS", 100),
			READ_BUFFER_SIZE:  getEnvAsInt("READ_BUFFER_SIZE", 1024),
			WRITE_BUFFER_SIZE: getEnvAsInt("WRITE_BUFFER_SIZE", 1024),
			PONG_WAIT:         getEnvAsInt("PONG_WAIT", 60),
			PING_PERIOD:       getEnvAsInt("PING_WAIT", 50),
			WRITE_WAIT:        getEnvAsInt("WRITE_WAIT", 10),
			MAX_MESSAGE:       getEnvAsInt("MAX_MESSAGE", 512),
		},
	}

	// Example: check for a critical setting
	if cfg.API_CONFIG.JWT_SECRET == "super-secret-key" {
		log.Println("WARNING: Using default JWT_SECRET. Update .env for production.")
	}

	return cfg, nil
}

func getEnvAsStr(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		result, err := strconv.Atoi(value)

		if err != nil {
			return result
		}
	}

	return fallback
}
