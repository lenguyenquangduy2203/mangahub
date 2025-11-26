package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	log.Println("WebSocket server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
