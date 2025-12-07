package main

import (
	"log"
	"mangahub/internal/tcp/server"
)

func main() {
	srv := server.NewProgressSyncServer("9001")
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
