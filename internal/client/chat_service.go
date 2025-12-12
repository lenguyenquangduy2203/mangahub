package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mangahub/pkg/models/dtos"
	"mangahub/pkg/utils"
	"mangahub/pkg/utils/colors"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func StartChat(token, room, username string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u, err := url.Parse(utils.WsURL)
	if err != nil {
		log.Fatalf("Invalid base url (%v)", err)
	}

	q := u.Query()
	q.Set("room", room)
	u.RawQuery = q.Encode()

	fmt.Print("\033[H\033[2J")
	fmt.Printf("%sWelcome to MangaHub Chat!%s\n", colors.ColorPurple, colors.ColorReset)
	fmt.Printf("Room: %s%s%s | User: %s%s%s\n", colors.ColorCyan, room, colors.ColorReset, colors.ColorGreen, username, colors.ColorReset)
	fmt.Println(strings.Repeat("-", 50))

	requestHeader := http.Header{}
	requestHeader.Add("Authorization", "Bearer "+strings.TrimSpace(token))

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), requestHeader)
	if err != nil {
		if resp != nil {
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(resp.Body)

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				log.Fatalf("Status %d: %s", resp.StatusCode, string(body))
			}
		}

		log.Fatalf("%sDial error: %v%s", colors.ColorRed, err, colors.ColorReset)
	}

	fmt.Println("\U000F0789 Connected! Type a message and press Enter (Ctrl+C to quit).")

	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Printf("%sError closing connection: %v%s", colors.ColorRed, err, colors.ColorReset)
		}
	}(c)

	done := make(chan struct{})

	// Reader Goroutine (Incoming Messages)
	go func() {
		defer close(done)
		for {
			_, rawMsg, err := c.ReadMessage()
			if err != nil {
				return
			}

			// Parse the JSON message
			var msg dtos.Message
			if err := json.Unmarshal(rawMsg, &msg); err != nil {
				fmt.Printf("\r%s[System] %s%s\n> ", colors.ColorYellow, string(rawMsg), colors.ColorReset)
				continue
			}

			displayTime := time.Now().Format("3:04PM, 02-Jan-06")
			if msg.Time != "" {
				if t, err := time.Parse(time.RFC3339, msg.Time); err == nil {
					displayTime = t.Format("3:04PM, 02-Jan-06")
				}
			}

			fmt.Print("\r\033[K")

			// Render based on message type
			switch msg.Type {
			case "chat":
				if msg.Username == username {
					fmt.Printf("[%s] %sYou:%s %s\n", displayTime, colors.ColorGreen, colors.ColorReset, msg.Text)
				} else {
					fmt.Printf("[%s] %s%s:%s %s\n", displayTime, colors.ColorCyan, msg.Username, colors.ColorReset, msg.Text)
				}
			case "join":
				fmt.Printf("%s\U000F0206 %s.%s\n", colors.ColorYellow, msg.Text, colors.ColorReset)
			case "leave":
				fmt.Printf("%s\U000F0A48 %s.%s\n", colors.ColorBlue, msg.Text, colors.ColorReset)
			default:
				// Generic fallback
				fmt.Printf("[%s] %s: %s\n", displayTime, msg.Username, msg.Text)
			}

			fmt.Print("> ")
		}
	}()

	// Input Scanner Goroutine (User Typing)
	inputChan := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			inputChan <- text
		}
	}()

	// Main Event Loop
	for {
		select {
		case <-done:
			fmt.Println("\nDisconnected from server.")
			return

		case text := <-inputChan:
			if text == "" {
				fmt.Print("\033[1A\033[K") // Move up and clear empty line
				fmt.Print("> ")
				continue
			}

			fmt.Print("\033[1A\033[K")

			// Send Message
			msg := dtos.Message{
				Type: "chat",
				Text: text,
			}
			jsonMsg, _ := json.Marshal(msg)

			err := c.WriteMessage(websocket.TextMessage, jsonMsg)
			if err != nil {
				log.Println("Write error:", err)
				return
			}

		case <-interrupt:
			fmt.Println("\n\U000F1821 Exiting chat...")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
