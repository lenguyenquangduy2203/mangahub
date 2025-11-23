package server

import "net"

type NotificationServer struct {
	Port    string
	Clients []net.UDPAddr
}
