package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ghostsecurity/reaper/internal/database/models"
)

func (h *Handler) GetRequests(c *fiber.Ctx) error {
	requests := []models.Request{}
	h.db.Preload("Response").Limit(250).Order("requests.created_at desc").Find(&requests)

	return c.JSON(requests)
}

func (h *Handler) GetRequest(c *fiber.Ctx) error {
	request := models.Request{}
	h.db.Preload("Response").First(&request, c.Params("id"))

	return c.JSON(request)
}
