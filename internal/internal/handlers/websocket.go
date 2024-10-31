package handlers

import (
	"log/slog"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	ws "github.com/ghostsecurity/reaper/internal/handlers/websocket"
)

func (h *Handler) WebSocketUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (h *Handler) WebSocket(c *websocket.Conn) {
	slog.Info("websocket client connected")

	client := &ws.Client{
		Conn: c,
		Pool: h.pool,
	}

	h.pool.Register <- client
	client.Read()
}
