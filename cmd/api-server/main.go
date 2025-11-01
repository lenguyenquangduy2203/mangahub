package main

import (
	"log"
	"net/http"
	"os"

	"mangahub/internal/api-server/handlers"
	"mangahub/internal/api-server/routes"
	"mangahub/internal/auth"
	"mangahub/internal/user"
	"mangahub/pkg/utils/config"
	"mangahub/pkg/utils/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Database Connection
	dbConnector, err := database.NewDatabaseConnection(cfg.DB_PATH)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Defer the close on the underlying SQL connection
	sqlConnection := dbConnector.(*database.Database).SQLDB()
	defer sqlConnection.Close()

	// Dependency Injection and Wiring
	// Repos & Utils
	userRepo := user.NewRepository(dbConnector)

	jwtSecret := cfg.JWT_SECRET
	jwtLifeTime := os.Getenv("JWT_ACCESS_TOKEN_LIFETIME")
	jwtManager, err := auth.NewJWTManager(jwtSecret, jwtLifeTime)
	if err != nil {
		log.Fatalf("Failed to initialize JWT Manager: %v", err)
	}
	passwordHasher := auth.NewPasswordUtility()

	// Services
	authService := auth.NewService(userRepo, jwtManager, passwordHasher)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Start server
	router := gin.Default()

	// version 1
	apiV1 := router.Group("/api/v1")

	{
		// Health endpoint
		apiV1.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "good",
			})
		})

		// Auth routes
		authGroupV1 := apiV1.Group("/auth")
		routes.AuthRoutesV1(authGroupV1, authHandler)
	}

	// Read port from env
	port := cfg.API_PORT

	router.Run(":" + port)
}
