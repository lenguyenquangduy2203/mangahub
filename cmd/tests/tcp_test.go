package tests

import (
	"bufio"
	"encoding/json"
	"net"
	"testing"
	"time"

	"mangahub/internal/tcp/models"
	"mangahub/internal/tcp/server"
)

func startTestServer(t *testing.T) (*server.ProgressSyncServer, string) {
	t.Helper()

	// 0 = dynamic free port
	server := server.NewProgressSyncServer("9000")

	go server.Start()
	time.Sleep(200 * time.Millisecond) // allow startup

	return server, server.Port
}

func connectClient(t *testing.T, port string) net.Conn {
	t.Helper()

	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		t.Fatalf("Failed to connect client: %v", err)
	}
	return conn
}

// ------------------------------------------------------------
// TEST 1 — Accepts multiple connections
// ------------------------------------------------------------
func TestServerAcceptsMultipleClients(t *testing.T) {
	server, port := startTestServer(t)
	defer server.Stop()

	c1 := connectClient(t, port)
	defer c1.Close()

	c2 := connectClient(t, port)
	defer c2.Close()

	time.Sleep(50 * time.Millisecond)
	if len(server.Connections) != 2 {
		t.Fatalf("Expected 2 connections, got %d", len(server.Connections))
	}
}

// ------------------------------------------------------------
// TEST 2 — Broadcast message to all clients
// ------------------------------------------------------------
func TestBroadcastProgressUpdate(t *testing.T) {
	server, port := startTestServer(t)
	defer server.Stop()

	c1 := connectClient(t, port)
	c2 := connectClient(t, port)
	defer c1.Close()
	defer c2.Close()

	msg := models.ProgressUpdate{
		UserID:    "u123",
		MangaID:   "onepiece",
		Chapter:   1090,
		Timestamp: time.Now().Unix(),
	}

	go func() {
		server.Broadcast <- msg
	}()

	// both clients must receive JSON
	received1 := readJSONMessage(t, c1)
	received2 := readJSONMessage(t, c2)

	if received1.UserID != msg.UserID ||
		received1.MangaID != msg.MangaID ||
		received1.Chapter != msg.Chapter {
		t.Fatal("Client 1 did not receive correct broadcast")
	}

	if received2.UserID != msg.UserID ||
		received2.MangaID != msg.MangaID ||
		received2.Chapter != msg.Chapter {
		t.Fatal("Client 2 did not receive correct broadcast")
	}
}

// ------------------------------------------------------------
// TEST 3 — Disconnect cleanup
// ------------------------------------------------------------
func TestClientDisconnectsCleanly(t *testing.T) {
	server, port := startTestServer(t)
	defer server.Stop()

	c := connectClient(t, port)

	time.Sleep(50 * time.Millisecond)
	if len(server.Connections) != 1 {
		t.Fatalf("Expected 1 connection, got %d", len(server.Connections))
	}

	c.Close()
	time.Sleep(200 * time.Millisecond)

	if len(server.Connections) != 0 {
		t.Fatalf("Expected 0 connections after disconnect, got %d", len(server.Connections))
	}
}

// ------------------------------------------------------------
// TEST 4 — Server handles multiple concurrent updates
// ------------------------------------------------------------
func TestConcurrentBroadcasts(t *testing.T) {
	server, port := startTestServer(t)
	defer server.Stop()

	c1 := connectClient(t, port)
	c2 := connectClient(t, port)
	defer c1.Close()
	defer c2.Close()

	for i := range 5 {
		msg := models.ProgressUpdate{UserID: "u1", MangaID: "onepiece", Chapter: 1000 + i}
		server.Broadcast <- msg

		// Immediately read from clients
		_ = readJSONMessage(t, c1)
		_ = readJSONMessage(t, c2)
	}
}

// ------------------------------------------------------------
// Helper: read JSON from TCP connection
// ------------------------------------------------------------
func readJSONMessage(t *testing.T, conn net.Conn) models.ProgressUpdate {
	t.Helper()

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	reader := bufio.NewReader(conn)

	line, err := reader.ReadBytes('\n')
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	var msg models.ProgressUpdate
	if err := json.Unmarshal(line, &msg); err != nil {
		t.Fatalf("Invalid JSON received: %v", err)
	}

	return msg
}
