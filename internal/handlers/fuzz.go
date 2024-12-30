package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ghostsecurity/reaper/internal/database/models"
	"github.com/ghostsecurity/reaper/internal/tools/fuzz"
)

type FuzzAttackType string

const (
	FuzzAttackTypeIDOR  FuzzAttackType = "idor"
	FuzzAttackTypeBrute FuzzAttackType = "brute"
)

type CreateFuzzAttackRequest struct {
	Type       FuzzAttackType `json:"type"`      // "idor" or "brute"
	Hostname   string         `json:"hostname"`   // Target hostname
	Param      string         `json:"param"`      // Parameter to test
	Dictionary []string       `json:"dictionary"` // For brute force attacks
	MaxSuccess int           `json:"maxSuccess"` // Maximum successful attempts
	MaxRPS     int           `json:"maxRPS"`     // Rate limiting
}

func (h *Handler) CreateFuzzAttack(c *fiber.Ctx) error {
	var req CreateFuzzAttackRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	// Create the attack record
	attack := &models.FuzzAttack{
		Type:   string(req.Type),
		Status: "pending",
	}
	if err := h.DB.Create(attack).Error; err != nil {
		return err
	}

	// Start the appropriate attack type
	go func() {
		switch req.Type {
		case FuzzAttackTypeBrute:
			err := fuzz.CreateBruteForceAttack(
				attack.ID,
				req.Hostname,
				req.Param,
				req.Dictionary,
				h.Websocket,
				h.DB,
				req.MaxSuccess,
				req.MaxRPS,
			)
			if err != nil {
				h.Logger.Error("Failed to run brute force attack", "error", err)
			}
		}
	}()

	return c.JSON(attack)
}
