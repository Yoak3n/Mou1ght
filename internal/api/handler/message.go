package handler

import (
	"Mou1ght/internal/api/controller"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
)

func CreateMessage(c *fiber.Ctx) error {
	req := &request.CreateMessageRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	err = controller.CreateMessage(req)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Create message successfully")
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
