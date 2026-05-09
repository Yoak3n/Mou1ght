package handler

import (
	"Mou1ght/internal/api/controller"
	"Mou1ght/internal/domain/entity"
	"Mou1ght/internal/domain/model/schema/request"
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

func CreateMessage(c *fiber.Ctx) error {
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

	err = controller.CreateMessage(req)
	if err != nil {
		log.Printf("CreateMessage controller error: %v\n", err)
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Create message successfully")
}

func VisitorID(c *fiber.Ctx) error {
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

func DeleteMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	err := controller.DeleteMessageByID(id)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Delete message successfully")
}

func UpdateMessage(c *fiber.Ctx) error {
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

	err = controller.UpdateMessage(req)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Update message successfully")
}

func UpdateMessagePosition(c *fiber.Ctx) error {
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

	err = controller.UpdateMessagePosition(req, false)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Update message position successfully")
}

func DetailMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	msg, err := controller.GetMessageByID(id)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, msg)
}

func ListMessage(c *fiber.Ctx) error {
	req := &request.MessageListRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	result := make(map[string]any)
	if req.DateRange == nil {
		msgs, err := controller.ListMessages(nil, req.Sort)
		if err == nil {
			result["messages"] = msgs
		}
	} else {
		msgs, err := controller.ListMessages(req.DateRange, req.Sort)
		if err == nil {
			result["messages"] = msgs
		}
	}
	return util.SuccessResponse(c, result)
}

func ListMessagePublic(c *fiber.Ctx) error {
	req := &request.MessageListRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	result := make(map[string]any)
	var msgs []*entity.MessageEntity
	if req.DateRange == nil {
		ms, err := controller.ListMessages(nil, req.Sort)
		if err == nil {
			msgs = ms
		}
	} else {
		ms, err := controller.ListMessages(req.DateRange, req.Sort)
		if err == nil {
			msgs = ms
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

func ViewMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	err := controller.ViewMessage(id)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil)
}

func LikeMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	err := controller.LikeMessage(id)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil)
}

func OwnedMessageIDs(c *fiber.Ctx) error {
	token := c.Query("token", "")
	if token == "" {
		return util.ErrorResponse(c, 400, "token is required")
	}
	jti, err := parseVisitorJTI(token)
	if err != nil {
		return util.ErrorResponse(c, 403, "Invalid visitor token")
	}
	ids, err := controller.GetOwnedMessageIDs(jti)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, fiber.Map{"ids": ids})
}
