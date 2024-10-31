package handlers

import (
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/ghostsecurity/reaper/internal/database/models"
	"github.com/ghostsecurity/reaper/internal/service"
	"github.com/ghostsecurity/reaper/internal/tools/fuzz"
	"github.com/ghostsecurity/reaper/internal/types"
)

func (h *Handler) GetEndpoints(c *fiber.Ctx) error {
	endpoints := []models.Endpoint{}
	h.db.Find(&endpoints)

	return c.JSON(endpoints)
}

func (h *Handler) GetEndpoint(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	endpoint := models.Endpoint{}
	h.db.First(&endpoint, id)

	return c.JSON(endpoint)
}

func (h *Handler) CreateEndpoint(c *fiber.Ctx) error {
	var endpointInput service.EndpointInput

	if err := c.BodyParser(&endpointInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if endpointInput.Hostname == "" || endpointInput.Path == "" || endpointInput.Method == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "hostname, path, and method are required"})
	}

	endpoint, err := service.CreateEndpoint(h.db, endpointInput)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(endpoint)
}

func (h *Handler) CreateAttack(c *fiber.Ctx) error {
	var atk struct {
		EndpointID      uint     `json:"endpoint_id"`
		Tags            []string `json:"tags"` // TODO: remove
		ExcludedKeys    []string `json:"excluded_keys"`
		ExcludedKeyType string   `json:"excluded_key_type"`
	}

	if err := c.BodyParser(&atk); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if atk.EndpointID < 1 || len(atk.Tags) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "endpoint_id and tags are required"})
	}

	// TODO: get domain from endpoint
	domain := "ghostbank.net"
	go func() {
		err := fuzz.CreateAttack(domain, atk.ExcludedKeys, h.pool, h.db, 100, 1000, 10)
		if err != nil {
			slog.Error("[create attack]", "msg", "error creating attack", "error", err)
		}
	}()

	return c.JSON(fiber.Map{"status": "ok"})
}

func (h *Handler) DeleteAttackResults(c *fiber.Ctx) error {
	// TODO: delete by endpoint id
	// id, err := strconv.Atoi(c.Params("id"))
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// }

	res := h.db.Delete(&models.FuzzResult{})
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": res.Error.Error()})
	}

	m := &types.AttackResultMessage{
		Type: types.MessageTypeAttackResultClear,
	}

	h.pool.Broadcast <- m

	return c.JSON(fiber.Map{"status": "ok", "deleted": res.RowsAffected})
}
