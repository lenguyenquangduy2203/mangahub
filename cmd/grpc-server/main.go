package main

import (
	"log"
	"mangahub/internal/api-server/mangas"
	"mangahub/internal/api-server/users"
	grpcsrv "mangahub/internal/grpc"
	"mangahub/internal/tcp/client"
	"mangahub/pkg/utils/config"
	"mangahub/pkg/utils/database"
	mangapb "mangahub/proto/manga"
	"net"

	"google.golang.org/grpc"
)

func main() {
	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Read port from env
	// port := cfg.PORT

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	progressClient := client.NewProgressSyncClient("mangahub-tcp:9000")
	if err := progressClient.Connect(); err != nil {
		log.Fatalf("TCP client error: %v", err)
	}

	// Initialize Database Connection
	dbConnector, err := database.NewDatabaseConnection(cfg.API_CONFIG.DB_PATH)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Defer the close on the underlying SQL connection
	sqlConnection := dbConnector.(*database.Database).SQLDB()
	defer sqlConnection.Close()

	// Dependency Injection and Wiring
	// Repos
	userRepo := users.NewRepository(dbConnector)
	mangaRepo := mangas.NewRepository(dbConnector)

	// Service
	userService := users.NewService(userRepo)
	mangaService := mangas.NewService(mangaRepo)
	mangaGrpcService := grpcsrv.NewMangaServiceGrpcServer(userService, mangaService, progressClient)

	mangapb.RegisterMangaServiceServer(s, mangaGrpcService)

	log.Println("gRPC server running at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	select {}
}
