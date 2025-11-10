package server

import (
	"fmt"
	"net"
	"sync"

	"mangahub/internal/tcp/models"
)

type ProgressSyncServer struct {
	Port        string
	Connections map[string]net.Conn
	Broadcast   chan models.ProgressUpdate
	mu          sync.Mutex
}

func NewProgressSyncServer(port string) *ProgressSyncServer {
	return &ProgressSyncServer{
		Port:        port,
		Connections: make(map[string]net.Conn),
		Broadcast:   make(chan models.ProgressUpdate, 10),
	}
}

func (s *ProgressSyncServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.Port))
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Printf("TCP Progress Sync Server started on port %s\n", s.Port)

	disconnect := make(chan string, 10)
	go startBroadcaster(s.Connections, s.Broadcast)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		addr := conn.RemoteAddr().String()
		s.mu.Lock()
		s.Connections[addr] = conn
		s.mu.Unlock()

		fmt.Printf("Client connected: %s\n", addr)
		go handleConnection(conn, s.Broadcast, disconnect)
	}
}
