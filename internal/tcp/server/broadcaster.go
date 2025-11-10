package server

import (
	"encoding/json"
	"net"

	"mangahub/internal/tcp/models"
)

func startBroadcaster(connections map[string]net.Conn, broadcast <-chan models.ProgressUpdate) {
	for update := range broadcast {
		data, _ := json.Marshal(update)
		data = append(data, '\n')

		for addr, conn := range connections {
			_, err := conn.Write(data)
			if err != nil {
				conn.Close()
				delete(connections, addr)
			}
		}
	}
}
