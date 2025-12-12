package main

import (
	"database/sql"
	"log"
	"mangahub/internal/socket"
	"net/http"

	"mangahub/internal/api-server/handlers"
	"mangahub/internal/api-server/mangas"
	"mangahub/internal/api-server/routes"
	"mangahub/internal/api-server/users"
	"mangahub/internal/auth"
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
	dbConnector, err := database.NewDatabaseConnection(cfg.API_CONFIG.DB_PATH)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Defer the close on the underlying SQL connection
	sqlConnection := dbConnector.(*database.Database).SQLDB()
	defer func(sqlConnection *sql.DB) {
		err := sqlConnection.Close()
		if err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}(sqlConnection)

	// Dependency Injection and Wiring
	// Repos
	userRepo := users.NewRepository(dbConnector)
	mangaRepo := mangas.NewRepository(dbConnector)

	// Utils
	jwtSecret := cfg.API_CONFIG.JWT_SECRET
	jwtLifeTime := cfg.API_CONFIG.JWT_ACCESS_TOKEN_LIFETIME
	jwtManager, err := auth.NewJWTManager(jwtSecret, jwtLifeTime)
	if err != nil {
		log.Fatalf("Failed to initialize JWT Manager: %v", err)
	}
	passwordHasher := auth.NewPasswordUtility()

	// Services
	authService := auth.NewService(userRepo, jwtManager, passwordHasher)
	userService := users.NewService(userRepo)
	mangaService := mangas.NewService(mangaRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	mangaHandler := handlers.NewMangaHandler(mangaService)

	// Middlewares
	jwtMiddleware := auth.AuthMiddleware(jwtManager)

	// Register socket hub
	hub := socket.NewHub(cfg.SOCKET_CONFIG)

	// Start server
	router := gin.Default()
	go hub.Run()

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

		// Manga routes
		mangaV1 := apiV1.Group("/manga")
		routes.MangaRoutesV1(mangaV1, mangaHandler)

		// User private routes
		userV1 := apiV1.Group("/users")
		userV1.Use(jwtMiddleware)
		routes.UserRoutesV1(userV1, userHandler)

		// Socket endpoint (private)
		apiV1.GET("/ws", jwtMiddleware, socket.HandleWebSocket(hub, cfg.SOCKET_CONFIG))
	}

	// Read port from env
	port := cfg.API_CONFIG.API_PORT

	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}
