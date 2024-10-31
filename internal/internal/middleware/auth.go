package middleware

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/ghostsecurity/reaper/internal/service"
)

const (
	AuthTokenHeader = "X-Reaper-Token"
)

func TokenAuth(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get(AuthTokenHeader)
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		user, err := service.GetUserByToken(token, db)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		c.Locals("user", user)
		return c.Next()
	}
}
