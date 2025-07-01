package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/ghostsecurity/reaper/internal/database/models"
	"github.com/ghostsecurity/reaper/internal/types"
)

func (h *Handler) GetReports(c *fiber.Ctx) error {
	reports := []models.Report{}
	h.db.Order("created_at DESC").Find(&reports)

	return c.JSON(reports)
}

func (h *Handler) GetReport(c *fiber.Ctx) error {
	report := models.Report{}
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "report id is required"})
	}

	if _, err := strconv.Atoi(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid report id"})
	}

	err := h.db.First(&report, id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(report)
}

func (h *Handler) CreateReport(c *fiber.Ctx) error {
	reportInput := struct {
		Domain   string `json:"domain"`
		Markdown string `json:"markdown"`
	}{}

	if err := c.BodyParser(&reportInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	report := models.Report{
		Domain:   reportInput.Domain,
		Markdown: reportInput.Markdown,
	}
	resp := h.db.Create(&report)
	if resp.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": resp.Error.Error()})
	}

	msg := &types.ReportMessage{
		Type:      types.MessageTypeReportStatus,
		Domain:    report.Domain,
		Markdown:  report.Markdown,
		Timestamp: report.CreatedAt,
	}
	h.pool.Broadcast <- msg

	return c.JSON(report)
}

func (h *Handler) DeleteReport(c *fiber.Ctx) error {
	report := models.Report{}
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "report id is required"})
	}

	if _, err := strconv.Atoi(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid report id"})
	}

	res := h.db.Delete(&report, id)
	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "report not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
