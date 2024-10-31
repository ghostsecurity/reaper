package websocket

import (
	"log/slog"

	"github.com/ghostsecurity/reaper/internal/types"
)

// BroadcastMessage is a message that can be broadcast to websocket clients
type BroadcastMessage interface {
	GetType() types.MessageType
}

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan BroadcastMessage
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan BroadcastMessage),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			slog.Info("websocket client added to pool", "size", len(pool.Clients))
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			slog.Info("websocket client disconnected from pool", "size", len(pool.Clients))
		case message := <-pool.Broadcast:
			for client := range pool.Clients {
				client.Mutex.Lock()

				if err := client.Conn.WriteJSON(message); err != nil {
					slog.Error("websocket error sending message to client", "error", err)
				}

				client.Mutex.Unlock()
			}
		}
	}
}
