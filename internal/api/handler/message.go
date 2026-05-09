package handler

import (
	"Mou1ght/internal/api/controller"
	"Mou1ght/internal/domain/entity"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/pkg/util"
	"log"

	"github.com/gofiber/fiber/v2"
)

func CreateMessage(c *fiber.Ctx) error {
	req := &request.CreateMessageRequest{}
	err := c.BodyParser(req)
	if err != nil {
		log.Printf("CreateMessage BodyParser error: %v\n", err)
		return util.ErrorResponse(c, 400, err.Error())
	}
	log.Printf("CreateMessage ip=%s ua_len=%d content_len=%d pos=(%d,%d,%d) author_ip_len=%d\n",
		c.IP(),
		len(c.Get("User-Agent")),
		len(req.Content),
		req.Position.X, req.Position.Y, req.Position.Z,
		len(req.AuthorIP),
	)

	token, claims, perr := util.ParseVisitorToken(req.AuthorIP)
	if perr != nil || token == nil || !token.Valid {
		log.Printf("CreateMessage invalid visitor token: err=%v valid=%v\n", perr, token != nil && token.Valid)
		return util.ErrorResponse(c, 403, "Invalid visitor token")
	}
	log.Printf("CreateMessage visitor token ok: jti=%s claim_ip=%s claim_ua_len=%d\n", claims.ID, claims.IP, len(claims.UA))
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
		log.Printf("UpdateMessagePosition BodyParser error: %v\n", err)
		return util.ErrorResponse(c, 400, err.Error())
	}
	log.Printf("UpdateMessagePosition ip=%s msg_id=%s pos=(%d,%d,%d) author_ip_len=%d\n",
		c.IP(),
		req.ID,
		req.Position.X, req.Position.Y, req.Position.Z,
		len(req.AuthorIP),
	)
	token, claims, perr := util.ParseVisitorToken(req.AuthorIP)
	if perr != nil || token == nil || !token.Valid {
		log.Printf("UpdateMessagePosition invalid visitor token: err=%v valid=%v\n", perr, token != nil && token.Valid)
		return util.ErrorResponse(c, 403, "Invalid visitor token")
	}
	log.Printf("UpdateMessagePosition visitor token ok: jti=%s claim_ip=%s\n", claims.ID, claims.IP)
	// Check for admin/user token to allow bypass of IP check
	// isAdmin := false
	// headers := c.GetReqHeaders()
	// tokenHeader, ok := headers["Authorization"]
	// if ok && len(tokenHeader) > 0 && len(tokenHeader[0]) > 7 {
	// 	tokenString := tokenHeader[0][7:]
	// 	token, _, err := util.ParseToken(tokenString)
	// 	if err == nil && token.Valid {
	// 		isAdmin = true
	// 	}
	// }

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
