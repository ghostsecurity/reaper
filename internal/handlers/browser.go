package handlers

import (
	"encoding/json"
	"log/slog"
	"net"

	"github.com/ghostsecurity/reaper/internal/browser"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type BrowserHandler struct {
	browser *browser.Browser
}

func NewBrowserHandler(b *browser.Browser) *BrowserHandler {
	return &BrowserHandler{
		browser: b,
	}
}

// HandleVNC handles the WebSocket connection for VNC
func (h *BrowserHandler) HandleVNC(c *websocket.Conn) {
	// Create a connection to the VNC server
	vncConn, err := net.Dial("tcp", "localhost:5900")
	if err != nil {
		slog.Error("Failed to connect to VNC server", "error", err)
		return
	}
	defer vncConn.Close()

	// Bidirectional copy between WebSocket and VNC connection
	go func() {
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				return
			}
			if mt == websocket.BinaryMessage {
				_, err = vncConn.Write(msg)
				if err != nil {
					return
				}
			}
		}
	}()

	buf := make([]byte, 4096)
	for {
		n, err := vncConn.Read(buf)
		if err != nil {
			return
		}
		err = c.WriteMessage(websocket.BinaryMessage, buf[:n])
		if err != nil {
			return
		}
	}
}

// HandleNavigate handles browser navigation requests
func (h *BrowserHandler) HandleNavigate(c *fiber.Ctx) error {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.browser.Navigate(req.URL); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// HandleReload handles browser reload requests
func (h *BrowserHandler) HandleReload(c *fiber.Ctx) error {
	if err := h.browser.Reload(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// HandleStart starts the browser instance
func (h *BrowserHandler) HandleStart(c *fiber.Ctx) error {
	if err := h.browser.Start(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// HandleStop stops the browser instance
func (h *BrowserHandler) HandleStop(c *fiber.Ctx) error {
	if err := h.browser.Stop(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
