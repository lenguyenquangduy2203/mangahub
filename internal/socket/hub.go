package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"mangahub/pkg/models/dtos"
	"mangahub/pkg/utils/config"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID       string
	Username string
	Conn     *websocket.Conn
	Room     string
	Send     chan []byte
}

type Room struct {
	Name    string
	Clients map[string]*Client
	History []dtos.Message
	Mu      sync.RWMutex
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Mu         sync.RWMutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewHub(cfg config.SocketConfig) *Hub {
	rooms := cfg.CHAT_ROOMS
	hub := &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}

	for i := range rooms {
		rooms[i] = strings.TrimSpace(rooms[i])
		if rooms[i] == "" {
			continue
		}

		hub.Rooms[rooms[i]] = &Room{
			Name:    rooms[i],
			Clients: make(map[string]*Client),
			History: []dtos.Message{},
		}
		log.Println("Initialized room:", rooms[i])
	}

	return hub
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.addClientToRoom(client)

		case client := <-h.Unregister:
			h.removeClientFromRoom(client)
		}
	}
}

func (h *Hub) addClientToRoom(client *Client) {
	h.Mu.RLock()
	room := h.Rooms[client.Room]
	h.Mu.RUnlock()

	if room == nil {
		log.Printf("BUG: room %s missing even though handler allowed join", client.Room)
		_ = client.Conn.Close()
		return
	}

	room.Mu.Lock()
	room.Clients[client.ID] = client
	room.Mu.Unlock()

	log.Printf("Client %s joined room %s", client.ID, client.Room)

	sendRoomHistory(room, client)

	msg := dtos.Message{
		Type: "join",
		Room: client.Room,
		Text: fmt.Sprintf("%s joined the room", client.Username),
		Time: time.Now().Format(time.RFC3339),
	}
	h.broadcastToRoom(client.Room, msg)
}

func sendRoomHistory(room *Room, client *Client) {
	room.Mu.RLock()
	history := make([]dtos.Message, len(room.History))
	copy(history, room.History)
	room.Mu.RUnlock()

	for _, msg := range history {
		data, _ := json.Marshal(msg)
		client.Send <- data
	}
}

func (h *Hub) removeClientFromRoom(client *Client) {
	h.Mu.RLock()
	room, exists := h.Rooms[client.Room]
	h.Mu.RUnlock()

	if !exists {
		return
	}

	room.Mu.Lock()
	if _, ok := room.Clients[client.ID]; ok {
		delete(room.Clients, client.ID)
		close(client.Send)
	}
	room.Mu.Unlock()

	log.Printf("Client %s left room %s (Remaining: %d)",
		client.Username, client.Room, len(room.Clients))

	// Send leave message to room
	msg := dtos.Message{
		Type: "leave",
		Room: client.Room,
		Text: fmt.Sprintf("%s left the room", client.Username),
		Time: time.Now().Format(time.RFC3339),
	}

	h.broadcastToRoom(client.Room, msg)
}

func (h *Hub) broadcastToRoom(roomName string, msg dtos.Message) {
	h.Mu.RLock()
	room := h.Rooms[roomName]
	h.Mu.RUnlock()

	if room == nil {
		return
	}

	room.Mu.Lock()
	if len(room.History) > 200 {
		room.History = room.History[len(room.History)-200:]
	}
	if msg.Type == "chat" {
		room.History = append(room.History, msg)
	}
	room.Mu.Unlock()

	data, _ := json.Marshal(msg)
	var dead []string

	room.Mu.RLock()
	for id, client := range room.Clients {
		select {
		case client.Send <- data:
		default:
			dead = append(dead, id)
		}
	}
	room.Mu.RUnlock()

	if len(dead) > 0 {
		room.Mu.Lock()
		for _, id := range dead {
			close(room.Clients[id].Send)
			delete(room.Clients, id)
		}
		room.Mu.Unlock()
	}
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	err := c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	if err != nil {
		return
	}
	c.Conn.SetPongHandler(func(string) error {
		err := c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			return err
		}
		return nil
	})

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg dtos.Message
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		// Set message metadata
		msg.Username = c.Username
		msg.Room = c.Room
		msg.Type = "chat"
		msg.Time = time.Now().Format(time.RFC3339)

		// Broadcast to room
		hub.broadcastToRoom(c.Room, msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		select {
		case message, ok := <-c.Send:
			err := c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				return
			}
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					return
				}
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			err := c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				return
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
