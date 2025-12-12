package socket

import (
	"fmt"
	"log"
	"mangahub/internal/auth"
	"mangahub/pkg/utils/config"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleWebSocket(hub *Hub, cfg config.SocketConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsRaw, exists := c.Get(auth.USER_KEY)
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		claims := claimsRaw.(*auth.JWTClaims)

		username := claims.UserName
		userId := claims.UserID

		roomName := c.Query("room")

		if username == "" || roomName == "" {
			c.JSON(400, gin.H{"error": "username and room required"})
			return
		}

		hub.Mu.RLock()
		room, exists := hub.Rooms[roomName]
		hub.Mu.RUnlock()

		if exists {
			room.Mu.RLock()
			_, userExists := room.Clients[userId]
			room.Mu.RUnlock()

			if userExists {
				c.JSON(409, gin.H{"error": fmt.Sprintf("YOU! Yes, YOU! are already in this '%s' room", roomName)})
				return
			}
		} else {
			log.Printf("Room %s does not exist. Rejecting client %s", roomName, userId)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Room '%s' does not exist, default rooms: [%s]", roomName, strings.Join(cfg.CHAT_ROOMS, ", "))})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Upgrade failed: %v", err)
			return
		}

		client := &Client{
			ID:       userId,
			Username: username,
			Room:     roomName,
			Conn:     conn,
			Send:     make(chan []byte, 256),
		}

		hub.Register <- client

		go client.writePump()
		go client.readPump(hub)
	}
}
