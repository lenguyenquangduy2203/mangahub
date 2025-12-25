package main

import (
	"log"
	"mangahub/internal/tcp/server"
)

func main() {
	srv := server.NewProgressSyncServer("9000")
	if err := srv.Start(); err != nil {
		log.Fatalf("TCP server error: %v", err)
	}

	select {}
}
