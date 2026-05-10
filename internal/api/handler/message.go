package handler

import (
	"Mou1ght/internal/api/service"
	"Mou1ght/internal/domain/entity"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"log"

	"github.com/gofiber/fiber/v2"
)

func parseVisitorJTI(authorIP string) (string, error) {
	token, claims, perr := util.ParseVisitorToken(authorIP)
	if perr != nil || token == nil || !token.Valid {
		return "", fiber.ErrForbidden
	}
	return claims.ID, nil
}

type MessageHandler struct {
	messageService *service.MessageService
	dtoService     *service.DTOService
}

func NewMessageHandler(messageService *service.MessageService, dtoService *service.DTOService) *MessageHandler {
	return &MessageHandler{messageService: messageService, dtoService: dtoService}
}

func (h *MessageHandler) CreateMessage(c *fiber.Ctx) error {
	req := &request.CreateMessageRequest{}
	err := c.BodyParser(req)
	if err != nil {
		log.Printf("CreateMessage BodyParser error: %v\n", err)
		return util.ErrorResponse(c, 400, err.Error())
	}

	jti, err := parseVisitorJTI(req.AuthorIP)
	if err != nil {
		return util.ErrorResponse(c, 403, "Invalid visitor token")
	}
	req.AuthorIP = jti

	err = h.messageService.CreateMessage(req)
	if err != nil {
		log.Printf("CreateMessage controller error: %v\n", err)
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Create message successfully")
}

func (h *MessageHandler) VisitorID(c *fiber.Ctx) error {
	ip := c.IP()
	ua := c.Get("User-Agent")
	id, err := util.ReleaseVisitorToken(ip, ua)
	if err != nil {
		log.Printf("VisitorID ReleaseVisitorToken error: %v\n", err)
		return util.ErrorResponse(c, 500, err.Error())
	}
	log.Printf("VisitorID issued: ip=%s ua_len=%d token_len=%d\n", ip, len(ua), len(id))
	return util.SuccessResponse(c, fiber.Map{
		"id": id,
	})
}

func (h *MessageHandler) DeleteMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	err := h.messageService.DeleteMessageByID(id)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Delete message successfully")
}

func (h *MessageHandler) UpdateMessage(c *fiber.Ctx) error {
	req := &request.UpdateMessageRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}

	jti, err := parseVisitorJTI(req.AuthorIP)
	if err != nil {
		return util.ErrorResponse(c, 403, "Invalid visitor token")
	}
	req.AuthorIP = jti

	err = h.messageService.UpdateMessage(req)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Update message successfully")
}

func (h *MessageHandler) UpdateMessagePosition(c *fiber.Ctx) error {
	req := &request.UpdateMessagePositionRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}

	jti, err := parseVisitorJTI(req.AuthorIP)
	if err != nil {
		return util.ErrorResponse(c, 403, "Invalid visitor token")
	}
	req.AuthorIP = jti

	err = h.messageService.UpdateMessagePosition(req, false)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Update message position successfully")
}

func (h *MessageHandler) DetailMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	msg, err := h.messageService.GetMessageByID(id)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	d := h.dtoService.GetMessageEntityFromTable(msg)
	return util.SuccessResponse(c, d)
}

func (h *MessageHandler) messages(msgs []*table.MessageTable) map[string]any {
	ret := make(map[string]any)
	mes := h.dtoService.GetMessagesEntityFromTables(msgs)
	ret["messages"] = mes
	return ret
}

func (h *MessageHandler) ListMessage(c *fiber.Ctx) error {
	req := &request.MessageListRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	var result map[string]any
	if req.DateRange == nil {
		msgs, err := h.messageService.ListMessages(nil, req.Sort)
		if err == nil {
			result = h.messages(msgs)
		}
	} else {
		msgs, err := h.messageService.ListMessages(req.DateRange, req.Sort)
		if err == nil {
			result = h.messages(msgs)
		}
	}
	return util.SuccessResponse(c, result)
}

func (h *MessageHandler) ListMessagePublic(c *fiber.Ctx) error {
	req := &request.MessageListRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	var result map[string]any
	var msgs []*entity.MessageEntity
	if req.DateRange == nil {
		ms, err := h.messageService.ListMessages(nil, req.Sort)
		if err == nil {
			msgs = h.dtoService.GetMessagesEntityFromTables(ms)
		}
	} else {
		ms, err := h.messageService.ListMessages(req.DateRange, req.Sort)
		if err == nil {
			msgs = h.dtoService.GetMessagesEntityFromTables(ms)
		}
	}
	filtered := make([]*entity.MessageEntity, 0, len(msgs))
	for _, m := range msgs {
		if m != nil && m.State.Status == "published" {
			filtered = append(filtered, m)
		}
	}
	result["messages"] = filtered
	return util.SuccessResponse(c, result)
}

func (h *MessageHandler) ViewMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	err := h.messageService.ViewMessage(id)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil)
}

func (h *MessageHandler) LikeMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	err := h.messageService.LikeMessage(id)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil)
}

func (h *MessageHandler) OwnedMessageIDs(c *fiber.Ctx) error {
	token := c.Query("token", "")
	if token == "" {
		return util.ErrorResponse(c, 400, "token is required")
	}
	jti, err := parseVisitorJTI(token)
	if err != nil {
		return util.ErrorResponse(c, 403, "Invalid visitor token")
	}
	ids, err := h.messageService.GetOwnedMessageIDs(jti)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, fiber.Map{"ids": ids})
}
