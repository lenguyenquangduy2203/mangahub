package server

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"time"

	"mangahub/internal/tcp/models"
)

func handleConnection(conn net.Conn, broadcast chan<- models.ProgressUpdate, disconnect chan<- string) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()

	reader := bufio.NewReader(conn)

	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				disconnect <- clientAddr
			}
			return
		}

		var update models.ProgressUpdate
		if err := json.Unmarshal(data, &update); err != nil {
			continue
		}

		update.Timestamp = time.Now().Unix()
		broadcast <- update
	}
}
