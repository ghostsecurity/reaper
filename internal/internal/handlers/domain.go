package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/ghostsecurity/reaper/internal/database"
	"github.com/ghostsecurity/reaper/internal/database/models"
	"github.com/ghostsecurity/reaper/internal/tools/scan"
	"github.com/ghostsecurity/reaper/internal/types"
)

func (h *Handler) GetDomains(c *fiber.Ctx) error {
	domains := []models.Domain{}
	h.db.Order("created_at desc").Find(&domains)

	return c.JSON(domains)
}

func (h *Handler) GetDomain(c *fiber.Ctx) error {
	domain := models.Domain{}
	res := h.db.First(&domain, c.Params("id"))
	if res.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": res.Error.Error()})
	}

	return c.JSON(domain)
}

func (h *Handler) CreateDomain(c *fiber.Ctx) error {
	domainInput := struct {
		Name     string `json:"name"`
		AutoScan bool   `json:"auto_scan"`
	}{}

	if err := c.BodyParser(&domainInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := validDomain(domainInput.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	domain := models.Domain{
		Name:   domainInput.Name,
		Status: types.DomainStatusPending,
	}
	resp := h.db.Create(&domain)
	if resp.Error != nil {
		if strings.HasPrefix(resp.Error.Error(), database.ErrUniqueConstraintFailed) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "domain already exists in this project"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": resp.Error.Error()})
	}

	if domainInput.AutoScan {
		go scan.FindSubdomains(domain, h.db, h.pool)
	}

	msg := &types.DomainMessage{
		Type:      types.MessageTypeScanDomain,
		Domain:    domain.Name,
		Timestamp: domain.CreatedAt,
	}
	h.pool.Broadcast <- msg

	return c.JSON(domain)
}

func (h *Handler) UpdateDomain(c *fiber.Ctx) error {
	domain := models.Domain{}
	if err := c.BodyParser(&domain); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	h.db.Model(&domain).Updates(domain)

	return c.JSON(domain)
}

func (h *Handler) DeleteDomain(c *fiber.Ctx) error {
	domain := models.Domain{}
	res := h.db.Delete(&domain, c.Params("id"))
	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "domain not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetDomainHosts(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	domain := models.Host{
		DomainID: uint(id),
	}
	hosts := []models.Host{}
	// TODO: pagination
	h.db.Where(domain).Limit(500).Find(&hosts)

	return c.JSON(hosts)
}

func validDomain(domain string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var r net.Resolver
	ns, err := r.LookupNS(ctx, domain)
	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("domain lookup timeout")
		}
		return fmt.Errorf("invalid domain, no nameserver found")
	}

	servers := make([]string, len(ns))
	for i, server := range ns {
		servers[i] = server.Host
	}

	slog.Info("scan.create", "domain", domain, "nameservers", strings.Join(servers, ","))

	return nil
}
