package handlers

import "github.com/gofiber/fiber/v2"

type FuzzAttack struct {
	Target string `json:"target"`
}

func (h *Handler) CreateFuzzAttack(c *fiber.Ctx) error {
	return nil
}
