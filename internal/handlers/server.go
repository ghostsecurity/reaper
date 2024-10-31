package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/ghostsecurity/reaper/internal/database/models"
	"github.com/ghostsecurity/reaper/internal/middleware"
	"github.com/ghostsecurity/reaper/internal/service"
	"github.com/ghostsecurity/reaper/internal/types"
)

type NavigationRecord struct {
	To   string `json:"to"`
	From string `json:"from"`
}

func (h *Handler) Status(c *fiber.Ctx) error {
	ip := c.IP()
	token := c.Get(middleware.AuthTokenHeader)

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"status": "unauthorized", "ip": ip})
	}

	user, _ := service.GetUserByToken(token, h.db)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"status": "no user", "ip": ip})
	}

	return c.JSON(fiber.Map{"status": "ok", "ip": ip, "user": user})
}

// Navigation handles navigation follow commands from the frontend.
func (h *Handler) Navigation(c *fiber.Ctx) error {
	// ignore navigation commands from non-leader
	user := c.Locals("user").(*models.User)

	if user.Role != types.UserRoleAdmin {
		return c.SendStatus(fiber.StatusNoContent)
	}

	var r NavigationRecord
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	h.pool.Broadcast <- &types.NavigationFollowMessage{
		Type: types.MessageTypeNavigationFollow,
		From: r.From,
		To:   r.To,
	}

	slog.Info("navigation", "message", "leader navigation", "to", r.To, "from", r.From)
	return c.JSON(fiber.Map{"navigation": "ok"})
}

func (h *Handler) Register(c *fiber.Ctx) error {
	input := struct {
		Username         string `json:"username"`
		RegistrationCode string `json:"invite_code"`
	}{}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	admin, _ := service.FirstAdmin(h.db)
	if admin == nil {
		// register an admin user
		u, err := service.CreateAdminUser(input.Username, h.db)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error creating admin user"})
		}

		return c.JSON(fiber.Map{"user": u})
	}

	// ***************************** //
	// TODO: remove after testing !! //
	// ***************************** //
	//
	// always login as admin
	user, _ := service.GetUserByToken(admin.Token, h.db)
	if user == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "admin user not found"})
	}
	return c.JSON(fiber.Map{"user": user})

	// get the guest token system setting
	// TODO: remove after testing
	// guestToken, _ := service.GetSettingByKey("guest_token", h.db)
	// if guestToken == nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "guest token is not set"})
	// }

	// slog.Info("register", "guest_token", *guestToken)

	// // guest users need an invite code
	// if input.RegistrationCode != *guestToken {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid registration code"})
	// }

	// // register a guest user
	// u, err := service.CreateGuestUser(input.Username, h.db)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error creating guest user"})
	// }

	// slog.Info("register", "username", input.Username, "code", input.RegistrationCode, "token", u.Token)

	// return c.JSON(fiber.Map{"user": u})
}
