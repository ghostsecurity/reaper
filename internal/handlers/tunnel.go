package handlers

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/ghostsecurity/reaper/internal/service"
	"github.com/ghostsecurity/reaper/internal/tools/tunnel"
)

const tunnelURLKey = "tunnel_url"

func (h *Handler) TunnelStart(c *fiber.Ctx) error {
	if h.tunnel != nil {
		return c.JSON(fiber.Map{"status": "ok", "message": "tunnel already running", "url": h.tunnel.URL})
	}
	h.tunnel = tunnel.NewTunnel()
	err := h.tunnel.Start()
	if err != nil {
		slog.Error("[tunnel start]", "msg", "error starting tunnel", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "error starting tunnel"})
	}

	return c.JSON(fiber.Map{"status": "ok", "url": h.tunnel.URL})
}

func (h *Handler) TunnelStop(c *fiber.Ctx) error {
	if h.tunnel != nil {
		h.tunnel.Stop()
		h.tunnel = nil
	}

	_ = service.DeleteSettingByKey(tunnelURLKey, h.db)

	return c.JSON(fiber.Map{"status": "ok"})
}

func (h *Handler) TunnelStatus(c *fiber.Ctx) error {
	if h.tunnel == nil {
		return c.Status(fiber.StatusGone).JSON(fiber.Map{"status": "error", "message": "tunnel not running"})
	}

	// tunnel.URL is nil if the tunnel had a problem starting
	if h.tunnel.URL == nil {
		if h.tunnel.Failed {
			h.tunnel = nil
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "tunnel had a problem"})
		}
	}

	guestToken, _ := service.GetSettingByKey("guest_token", h.db)

	url := fmt.Sprintf("%s/?code=%s", *h.tunnel.URL, *guestToken)

	return c.JSON(fiber.Map{"status": "ok", "url": url})
}
