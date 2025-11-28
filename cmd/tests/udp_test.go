package tests

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	"mangahub/internal/udp/models"
	"mangahub/internal/udp/server"
)

// ------------------------------------------------------------
// Helper: Start the UDP server
// ------------------------------------------------------------
func startUDPServer(t *testing.T) (*server.NotificationServer, int) {
	t.Helper()

	// Let OS choose an available UDP port
	addr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		t.Fatalf("ResolveUDPAddr failed: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatalf("ListenUDP failed: %v", err)
	}

	_, portStr, _ := net.SplitHostPort(conn.LocalAddr().String())
	conn.Close()

	srv := &server.NotificationServer{Port: portStr}
	go srv.Listen()

	// give server time to start
	time.Sleep(30 * time.Millisecond)

	port, _ := strconv.Atoi(portStr)
	return srv, port
}

// ------------------------------------------------------------
// Helper: Create UDP client
// ------------------------------------------------------------
func newUDPClient(t *testing.T, port int) (*net.UDPConn, *net.UDPAddr) {
	t.Helper()

	serverAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(port))
	local := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	conn, err := net.ListenUDP("udp", local)
	if err != nil {
		t.Fatalf("ListenUDP client failed: %v", err)
	}
	return conn, serverAddr
}

// ------------------------------------------------------------
// Helper: Read packet with timeout
// ------------------------------------------------------------
func readPacket(t *testing.T, conn *net.UDPConn) []byte {
	t.Helper()

	buf := make([]byte, 2048)
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		t.Fatalf("Failed to read UDP packet: %v", err)
	}
	return buf[:n]
}

// ------------------------------------------------------------
// TEST 1 — Client Registration
// ------------------------------------------------------------
func TestClientRegistration(t *testing.T) {
	srv, port := startUDPServer(t)

	client, serverAddr := newUDPClient(t, port)
	defer client.Close()

	client.WriteToUDP([]byte("register"), serverAddr)
	time.Sleep(50 * time.Millisecond)

	if len(srv.Clients) != 1 {
		t.Fatalf("expected 1 registered client, got %d", len(srv.Clients))
	}
}

// ------------------------------------------------------------
// TEST 2 — Broadcast a notification to registered clients
// ------------------------------------------------------------
func TestBroadcastNotification(t *testing.T) {
	srv, port := startUDPServer(t)

	c1, serverAddr := newUDPClient(t, port)
	c2, _ := newUDPClient(t, port)
	defer c1.Close()
	defer c2.Close()

	// register both clients
	c1.WriteToUDP([]byte("register"), serverAddr)
	c2.WriteToUDP([]byte("register"), serverAddr)
	time.Sleep(50 * time.Millisecond)

	notification := models.Notification{
		Type:      "chapter",
		MangaID:   "op",
		Message:   "Chapter 100 released",
		Timestamp: time.Now().Unix(),
	}

	srv.Broadcast(notification)

	p1 := readPacket(t, c1)
	p2 := readPacket(t, c2)

	var n1, n2 models.Notification
	json.Unmarshal(p1, &n1)
	json.Unmarshal(p2, &n2)

	if n1.MangaID != "op" || n2.MangaID != "op" {
		t.Fatalf("broadcast content mismatch")
	}
}

// ------------------------------------------------------------
// TEST 3 — Duplicate registration not added twice
// ------------------------------------------------------------
func TestDuplicateRegistration(t *testing.T) {
	srv, port := startUDPServer(t)

	client, addr := newUDPClient(t, port)
	defer client.Close()

	client.WriteToUDP([]byte("register"), addr)
	client.WriteToUDP([]byte("register"), addr)
	time.Sleep(50 * time.Millisecond)

	if len(srv.Clients) != 1 {
		t.Fatalf("expected 1 unique client, got %d", len(srv.Clients))
	}
}

// ------------------------------------------------------------
// TEST 4 — Broadcasting with no clients should not crash
// ------------------------------------------------------------
func TestBroadcastNoClients(t *testing.T) {
	srv, _ := startUDPServer(t)

	srv.Broadcast(models.Notification{
		Type:    "info",
		Message: "No clients test",
	})

	// pass if no panic
}

// ------------------------------------------------------------
// TEST 5 — Client receives multiple notifications
// ------------------------------------------------------------
func TestMultipleNotifications(t *testing.T) {
	srv, port := startUDPServer(t)

	c, addr := newUDPClient(t, port)
	defer c.Close()

	c.WriteToUDP([]byte("register"), addr)
	time.Sleep(50 * time.Millisecond)

	for i := 1; i <= 3; i++ {
		srv.Broadcast(models.Notification{
			Type:    "update",
			Message: fmt.Sprintf("msg-%d", i),
		})
	}

	count := 0
	for i := 0; i < 3; i++ {
		readPacket(t, c)
		count++
	}

	if count != 3 {
		t.Fatalf("expected 3 notifications, got %d", count)
	}
}
