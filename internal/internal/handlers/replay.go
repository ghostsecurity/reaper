package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/ghostsecurity/reaper/internal/tools/replay"
	"github.com/ghostsecurity/reaper/internal/types"
)

func (h *Handler) Replay(c *fiber.Ctx) error {
	replayInput := types.ReplayInput{}

	if err := c.BodyParser(&replayInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := replay.Do(context.TODO(), &replayInput)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}
