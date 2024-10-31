package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/ghostsecurity/reaper/internal/types"
)

// unused
func (h *Handler) ExploreHostExample(c *fiber.Ctx) error {
	msg := types.ExploreHostMessage{
		Type:      "explore_host",
		Name:      "example.com",
		Timestamp: time.Now(),
	}
	h.pool.Broadcast <- &msg

	return c.JSON(fiber.Map{"explore host": "ok"})
}

// unused
func (h *Handler) ExploreEndpointExample(c *fiber.Ctx) error {
	msg := types.ExploreEndpointMessage{
		Type:      "explore_endpoint",
		Path:      "/api/v2/users",
		Host:      "example.com",
		Timestamp: time.Now(),
	}
	h.pool.Broadcast <- &msg

	return c.JSON(fiber.Map{"explore endpoint": "ok"})
}
