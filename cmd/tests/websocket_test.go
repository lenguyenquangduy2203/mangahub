package tests

import (
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"mangahub/internal/auth"
	"mangahub/internal/socket"
	"mangahub/pkg/models/dtos"
	"mangahub/pkg/utils/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

//
// -------------------- Test Auth Middleware --------------------
//

func testAuthMiddleware(userID, username string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(auth.USER_KEY, &auth.JWTClaims{
			UserID:   userID,
			UserName: username,
		})
		c.Next()
	}
}

//
// -------------------- Test Server Setup --------------------
//

func startTestWSServer(t *testing.T) (*socket.Hub, *httptest.Server) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	cfg := config.SocketConfig{
		CHAT_ROOMS: []string{"general"},
	}

	hub := socket.NewHub(cfg)
	go hub.Run()

	r := gin.New()
	r.Use(testAuthMiddleware("user-1", "tester"))
	r.GET("/ws", socket.HandleWebSocket(hub, cfg))

	srv := httptest.NewServer(r)
	return hub, srv
}

//
// -------------------- WebSocket Helper --------------------
//

func connectWS(t *testing.T, serverURL string) *websocket.Conn {
	t.Helper()

	u, _ := url.Parse(serverURL)
	u.Scheme = "ws"
	u.Path = "/ws"
	q := u.Query()
	q.Set("room", "general")
	u.RawQuery = q.Encode()

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	require.NoError(t, err)

	return ws
}

//
// -------------------- Tests --------------------
//

func TestWebSocket_JoinBroadcast(t *testing.T) {
	hub, srv := startTestWSServer(t)
	defer srv.Close()

	ws := connectWS(t, srv.URL)
	defer ws.Close()

	// First message should be "join"
	_, data, err := ws.ReadMessage()
	require.NoError(t, err)

	var msg dtos.Message
	err = json.Unmarshal(data, &msg)
	require.NoError(t, err)

	require.Equal(t, "join", msg.Type)
	require.Equal(t, "general", msg.Room)

	// Ensure client registered
	time.Sleep(50 * time.Millisecond)

	hub.Mu.RLock()
	room := hub.Rooms["general"]
	hub.Mu.RUnlock()

	room.Mu.RLock()
	require.Len(t, room.Clients, 1)
	room.Mu.RUnlock()
}

func TestWebSocket_ChatBroadcast(t *testing.T) {
	_, srv := startTestWSServer(t)
	defer srv.Close()

	ws := connectWS(t, srv.URL)
	defer ws.Close()

	// Drain join message
	ws.ReadMessage()

	// Send chat message
	input := dtos.Message{
		Text: "hello world",
	}
	err := ws.WriteJSON(input)
	require.NoError(t, err)

	// Read broadcast
	_, data, err := ws.ReadMessage()
	require.NoError(t, err)

	var msg dtos.Message
	err = json.Unmarshal(data, &msg)
	require.NoError(t, err)

	require.Equal(t, "chat", msg.Type)
	require.Equal(t, "hello world", msg.Text)
	require.Equal(t, "general", msg.Room)
	require.Equal(t, "tester", msg.Username)
}

func TestWebSocket_DisconnectCleanup(t *testing.T) {
	hub, srv := startTestWSServer(t)
	defer srv.Close()

	ws := connectWS(t, srv.URL)
	ws.ReadMessage() // join
	ws.Close()

	time.Sleep(100 * time.Millisecond)

	hub.Mu.RLock()
	room := hub.Rooms["general"]
	hub.Mu.RUnlock()

	room.Mu.RLock()
	require.Len(t, room.Clients, 0)
	room.Mu.RUnlock()
}
