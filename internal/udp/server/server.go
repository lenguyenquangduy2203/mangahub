package server

import (
	"fmt"
	"log"
	"net"
)

type NotificationServer struct {
	Port    string
	Clients []net.UDPAddr
}

func (s *NotificationServer) Listen() error {
	addr, err := net.ResolveUDPAddr("udp", ":"+s.Port)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Printf("UDP Notification Broadcast Server started on port %s\n", s.Port)

	buf := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Read error:", err)
			continue
		}

		msg := string(buf[:n])
		if msg == "register" {
			s.addClient(*clientAddr)
			log.Println("Registered client:", clientAddr)
		}
	}
}

func (s *NotificationServer) addClient(client net.UDPAddr) {
	for _, c := range s.Clients {
		if c.IP.Equal(client.IP) && c.Port == client.Port {
			return // already registered
		}
	}
	s.Clients = append(s.Clients, client)
}
