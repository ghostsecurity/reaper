package websocket

import (
	"log"
	"log/slog"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	ID    string
	Conn  *websocket.Conn
	Pool  *Pool
	Mutex sync.Mutex
}

type Message struct {
	Type int         `json:"type"`
	Body interface{} `json:"body"`
}

// Read listens for new messages being sent to the websocket connection
func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		msg := Message{
			Type: messageType,
			Body: string(p),
		}

		// if message body is "ping" then respond with "pong"
		// this manual heartbeat sent from the browser helps keep
		// the connection alive because browsers don't have a
		// native WebSocket ping
		if msg.Body == "ping" {
			c.Conn.WriteMessage(websocket.TextMessage, []byte("pong"))
		} else {
			// we don't need to log pings, but we do want to log everything else
			slog.Info("[ws]", "body", msg.Body)
		}
	}
}
