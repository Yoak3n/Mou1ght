package handler

import (
	"Mou1ght/internal/api/service"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"errors"
	"log"
	"strings"
	"time"

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

	jti, err := parseVisitorJTI(req.VisitorToken)
	if err != nil {
		return util.ErrorResponse(c, 403, "Invalid visitor token")
	}
	req.VisitorToken = jti

	err = h.messageService.CreateMessage(req)
	if err != nil {
		log.Printf("CreateMessage controller error: %v\n", err)
		ferr := &fiber.Error{}
		if errors.As(err, &ferr) {
			return util.ErrorResponse(c, ferr.Code, ferr.Message)
		}
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Create message successfully")
}

func (h *MessageHandler) VisitorID(c *fiber.Ctx) error {
	// 游客 token 铸造限流：浏览器经 server action 转发时共享同一来源 IP，
	// 限流上限按全局并发访问量设置，防批量刷身份。
	if !util.Allow("visitor:mint:"+util.ClientIP(c), 120, time.Minute) {
		return util.ErrorResponse(c, 429, "too many requests")
	}
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

	jti, err := parseVisitorJTI(req.VisitorToken)
	if err != nil {
		return util.ErrorResponse(c, 403, "Invalid visitor token")
	}
	req.VisitorToken = jti

	err = h.messageService.UpdateMessage(req)
	if err != nil {
		ferr := &fiber.Error{}
		if errors.As(err, &ferr) {
			return util.ErrorResponse(c, ferr.Code, ferr.Message)
		}
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Update message successfully")
}

// UpdateOwnMessage 访客编辑自己的留言：身份经游客 token 校验，服务端再核对所有权。
func (h *MessageHandler) UpdateOwnMessage(c *fiber.Ctx) error {
	req := &request.UpdateMessageRequest{}
	if err := c.BodyParser(req); err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}

	jti, err := parseVisitorJTI(req.VisitorToken)
	if err != nil {
		return util.ErrorResponse(c, 403, "Invalid visitor token")
	}
	req.VisitorToken = jti

	if err := h.messageService.UpdateMessage(req); err != nil {
		ferr := &fiber.Error{}
		if errors.As(err, &ferr) {
			return util.ErrorResponse(c, ferr.Code, ferr.Message)
		}
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Update message successfully")
}

// DeleteOwnMessage 访客删除自己的留言。
func (h *MessageHandler) DeleteOwnMessage(c *fiber.Ctx) error {
	req := &request.DeleteOwnMessageRequest{}
	if err := c.BodyParser(req); err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}

	jti, err := parseVisitorJTI(req.VisitorToken)
	if err != nil {
		return util.ErrorResponse(c, 403, "Invalid visitor token")
	}

	if err := h.messageService.DeleteOwnMessage(req.ID, jti); err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Delete message successfully")
}

func (h *MessageHandler) UpdateMessagePosition(c *fiber.Ctx) error {
	req := &request.UpdateMessagePositionRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}

	jti, err := parseVisitorJTI(req.VisitorToken)
	if err != nil {
		return util.ErrorResponse(c, 403, "Invalid visitor token")
	}
	req.VisitorToken = jti

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
	msgs, total, err := h.messageService.ListMessages(req.DateRange, req.Sort, req.Page, req.PageSize)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	result := h.messages(msgs)
	result["total"] = total
	if req.Page > 0 && req.PageSize > 0 {
		result["page"] = req.Page
		result["page_size"] = req.PageSize
	}
	return util.SuccessResponse(c, result)
}

func (h *MessageHandler) ListMessagePublic(c *fiber.Ctx) error {
	req := &request.MessageListRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	ms, total, err := h.messageService.ListMessagesPublic(req.DateRange, req.Sort, req.Page, req.PageSize)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	result := h.messages(ms)
	result["total"] = total
	if req.Page > 0 && req.PageSize > 0 {
		result["page"] = req.Page
		result["page_size"] = req.PageSize
	}
	return util.SuccessResponse(c, result)
}

func (h *MessageHandler) ViewMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	identity, _, tokenRequired := util.VisitorIdentity(c)
	if tokenRequired {
		return util.ErrorResponse(c, 403, "visitor token is required")
	}
	if !util.Allow("view:cooldown:message:"+id+":"+identity, 1, time.Minute) {
		return util.SuccessResponse(c, nil)
	}
	if !util.Allow("view:flood:"+identity, 120, time.Minute) {
		return util.ErrorResponse(c, 429, "too many requests")
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
	identity, _, tokenRequired := util.VisitorIdentity(c)
	if tokenRequired {
		return util.ErrorResponse(c, 403, "visitor token is required")
	}
	if !util.Allow("like:dedup:message:"+id+":"+identity, 1, 24*time.Hour) {
		return util.SuccessResponse(c, nil)
	}
	if !util.Allow("like:flood:"+identity, 30, time.Minute) {
		return util.ErrorResponse(c, 429, "too many requests")
	}
	err := h.messageService.LikeMessage(id)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil)
}

func (h *MessageHandler) OwnedMessageIDs(c *fiber.Ctx) error {
	// 优先从 Authorization 头取游客 token（JWT 不进 URL），保留 query 参数兼容旧调用
	token := c.Query("token", "")
	if token == "" {
		if auth := c.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
			token = strings.TrimPrefix(auth, "Bearer ")
		}
	}
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
