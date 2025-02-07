package handlers

import (
	"log/slog"
	"strconv"
	"strings"

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

	endpoint, err := service.CreateOrUpdateEndpoint(h.db, endpointInput)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(endpoint)
}

func (h *Handler) CreateAttack(c *fiber.Ctx) error {
	var atk struct {
		EndpointID uint     `json:"endpoint_id"`
		Params     []string `json:"params"`
	}

	if err := c.BodyParser(&atk); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if atk.EndpointID < 1 || len(atk.Params) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "endpoint_id and params are required"})
	}

	// get hostname from endpoint
	endpoint := models.Endpoint{}
	err := h.db.First(&endpoint, atk.EndpointID).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "endpoint not found"})
	}

	attack := &models.FuzzAttack{
		Type:   "param",
		Params: strings.Join(atk.Params, ","),
	}
	h.db.Create(attack)

	go func() {
		err := fuzz.CreateAttack(attack.ID, endpoint.Hostname, atk.Params, h.pool, h.db, 100, 1000, 5)
		if err != nil {
			slog.Error("[create attack]", "msg", "error creating attack", "error", err)
		}
	}()

	return c.JSON(attack)
}

func (h *Handler) GetAttacks(c *fiber.Ctx) error {
	attacks := []models.FuzzAttack{}
	h.db.Find(&attacks)

	return c.JSON(attacks)
}

func (h *Handler) GetAttack(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	attack := models.FuzzAttack{}
	h.db.First(&attack, id)

	return c.JSON(attack)
}

func (h *Handler) GetAttackResults(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	limit, err := strconv.Atoi(c.Query("limit", "50"))
	if err != nil {
		limit = 50
	}

	results := []models.FuzzResult{}
	h.db.Where("fuzz_attack_id = ?", id).Limit(limit).Find(&results)

	return c.JSON(results)
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
