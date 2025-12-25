package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"mangahub/internal/tcp/models"
)

// ProgressSyncServer is a TCP server for broadcasting progress updates.
type ProgressSyncServer struct {
	Port        string
	Connections map[net.Conn]struct{}
	Broadcast   chan models.ProgressUpdate
	listener    net.Listener
	mu          sync.Mutex
}

// NewProgressSyncServer creates a new server instance
func NewProgressSyncServer(port string) *ProgressSyncServer {
	return &ProgressSyncServer{
		Port:        port,
		Connections: make(map[net.Conn]struct{}),
		Broadcast:   make(chan models.ProgressUpdate, 10),
	}
}

// Start begins accepting TCP connections
func (s *ProgressSyncServer) Start() error {
	ln, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		return err
	}
	s.listener = ln
	fmt.Printf("TCP Progress Sync Server started on port %s\n", s.Port)

	go s.acceptLoop()
	go s.broadcastLoop()

	return nil
}

// Stop closes all connections and the listener
func (s *ProgressSyncServer) Stop() {
	fmt.Println("Stopping TCP server...")
	s.listener.Close()

	s.mu.Lock()
	for conn := range s.Connections {
		conn.Close()
	}
	s.Connections = make(map[net.Conn]struct{})
	s.mu.Unlock()

	close(s.Broadcast)
	fmt.Println("TCP server stopped.")
}

// acceptLoop accepts new connections
func (s *ProgressSyncServer) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}

		s.mu.Lock()
		s.Connections[conn] = struct{}{}
		s.mu.Unlock()

		fmt.Printf("Client connected: %s\n", conn.RemoteAddr())
		go s.handleConnection(conn)
	}
}

// handleConnection reads messages and removes client on disconnect
func (s *ProgressSyncServer) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		s.mu.Lock()
		delete(s.Connections, conn)
		s.mu.Unlock()
		fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())
	}()

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return
		}

		var update models.ProgressUpdate
		if err := json.Unmarshal(line, &update); err != nil {
			continue
		}

		update.Timestamp = time.Now().Unix()
		s.Broadcast <- update
	}
}

// broadcastLoop sends updates to all clients
func (s *ProgressSyncServer) broadcastLoop() {
	for update := range s.Broadcast {
		data, _ := json.Marshal(update)
		data = append(data, '\n')

		s.mu.Lock()
		for conn := range s.Connections {
			// Non-blocking write with short timeout
			conn.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
			if _, err := conn.Write(data); err != nil {
				conn.Close()
				delete(s.Connections, conn)
				fmt.Printf("Client disconnected during broadcast: %s\n", conn.RemoteAddr())
			}
		}
		s.mu.Unlock()
	}
}
