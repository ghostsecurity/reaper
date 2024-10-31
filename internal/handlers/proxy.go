package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ghostsecurity/reaper/internal/tools/proxy"
)

func (h *Handler) ProxyStart(c *fiber.Ctx) error {
	if h.proxy != nil {
		return c.JSON(fiber.Map{"status": "ok", "message": "proxy already running"})
	}
	h.proxy = proxy.NewProxy(h.pool, h.db)
	h.proxy.Start()

	return c.JSON(fiber.Map{"status": "ok"})
}

func (h *Handler) ProxyStop(c *fiber.Ctx) error {
	if h.proxy != nil {
		h.proxy.Stop()
		h.proxy = nil
	}
	return c.JSON(fiber.Map{"status": "ok"})
}

func (h *Handler) ProxyStatus(c *fiber.Ctx) error {
	if h.proxy == nil {
		return c.Status(fiber.StatusGone).JSON(fiber.Map{"status": "error", "message": "proxy not running"})
	}
	return c.JSON(fiber.Map{"status": "ok", "message": "proxy running"})
}
