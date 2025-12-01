package main

import (
	"log"
	grpcsrv "mangahub/internal/grpc"
	"mangahub/internal/mangas"
	"mangahub/internal/users"
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

	// Initialize Database Connection
	dbConnector, err := database.NewDatabaseConnection(cfg.DB_PATH)
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
	mangaGrpcService := grpcsrv.NewMangaServiceGrpcServer(*userService, *mangaService)

	mangapb.RegisterMangaServiceServer(s, mangaGrpcService)

	log.Println("gRPC server running at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
