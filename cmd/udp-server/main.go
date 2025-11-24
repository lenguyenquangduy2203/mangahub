package main

import (
	"log"
	"mangahub/internal/udp/server"
)

func main() {
	server := &server.NotificationServer{Port: "9001"}
	if err := server.Listen(); err != nil {
		log.Fatalf("UDP server error: %v", err)
	}
}
