package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/ghostsecurity/reaper/internal/database/models"
	"github.com/ghostsecurity/reaper/internal/service"
	"github.com/ghostsecurity/reaper/internal/tools/agent"
	"github.com/ghostsecurity/reaper/internal/types"
)

func (h *Handler) GetSessions(c *fiber.Ctx) error {
	sessions := []models.AgentSession{}
	h.db.Preload("Messages").Order("created_at desc").Find(&sessions)

	return c.JSON(sessions)
}

func (h *Handler) GetSession(c *fiber.Ctx) error {
	session := models.AgentSession{}
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "session id is required"})
	}

	if _, err := strconv.Atoi(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid session id"})
	}
	res := h.db.First(&session, id)
	if res.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": res.Error.Error()})
	}

	return c.JSON(session)
}

func (h *Handler) CreateSession(c *fiber.Ctx) error {
	sessionInput := struct {
		Description string `json:"description"`
	}{}

	if err := c.BodyParser(&sessionInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	session := models.AgentSession{
		Description: sessionInput.Description,
	}
	resp := h.db.Create(&session)
	if resp.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": resp.Error.Error()})
	}

	return c.JSON(session)
}

func (h *Handler) DeleteSession(c *fiber.Ctx) error {
	session := models.AgentSession{}
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "session id is required"})
	}

	if _, err := strconv.Atoi(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid session id"})
	}
	res := h.db.Delete(&session, id)
	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "session not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetSessionMessages(c *fiber.Ctx) error {
	session := models.AgentSession{}

	s := c.Params("id")
	if s == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "session id is required"})
	}

	// TODO: preload author
	res := h.db.Preload("Messages").First(&session, s)
	if res.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": res.Error.Error()})
	}

	return c.JSON(session.Messages)
}

func (h *Handler) CreateSessionMessage(c *fiber.Ctx) error {
	author := c.Locals("user").(*models.User)

	messageInput := struct {
		Message string `json:"content"`
	}{}

	s, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid session id"})
	}
	if s == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "session id is required"})
	}

	if err := c.BodyParser(&messageInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if messageInput.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "message is required"})
	}

	message := models.AgentSessionMessage{
		AuthorID:       author.ID,
		AuthorRole:     author.Role,
		AgentSessionID: uint(s),
		Content:        messageInput.Message,
	}

	resp := h.db.Create(&message)
	if resp.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": resp.Error.Error()})
	}

	msg := &types.AgentSessionMessage{
		Type:       types.MessageTypeAgentSessionMessage,
		AuthorID:   message.AuthorID,
		AuthorRole: author.Role,
		SessionID:  uint(s),
		Content:    message.Content,
	}

	h.pool.Broadcast <- msg

	// don't start a workflow if the message is from the agent
	if author.Role == types.UserRoleAgent {
		return c.JSON(message)
	}

	// start the agent workflow
	agentToken, _ := service.GetAgentToken(h.db)
	prompt := message.Content

	fmt.Printf("[***] running agent session %d, agent token: %s with prompt %s\n", s, agentToken, prompt)

	agentMgr := &agent.AgentManager{
		Ctx:        c,
		Pool:       h.pool,
		DB:         h.db,
		SessionID:  uint(s),
		Author:     author.ID,
		AuthorRole: types.UserRoleAgent,
		Prompt:     prompt,
	}

	agentMgr.StartAgent()

	return c.JSON(message)
}
