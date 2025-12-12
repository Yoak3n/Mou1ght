package handler

import (
	"Mou1ght/internal/api/controller"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
)

func ListPost(c *fiber.Ctx) error {
	req := &request.PostListRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	resultMap := make(map[string]interface{})
	// 除all外暂时未支持date_range，看需求是否需要
	switch req.Filter.Typ {
	case "category":
		resultMap = controller.CategoryListWithArticle(req)
	case "tag":
		resultMap = controller.TagListWithPost(req, req.IsSharing)
	case "author":
		resultMap = controller.AuthorListWithPost(req)
	case "all":
		resultMap = controller.AllListWithPost(req)
	default:
		return util.ErrorResponse(c, 400, "Invalid filter type")
	}

	return util.SuccessResponse(c, resultMap)
}

func ViewPost(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	typ := c.Query("type", "article")
	switch typ {
	case "article":
		err := controller.ViewArticle(id)
		if err != nil {
			return util.ErrorResponse(c, 500, err.Error())
		}
	case "sharing":
		err := controller.ViewSharing(id)
		if err != nil {
			return util.ErrorResponse(c, 500, err.Error())
		}
	case "message":
		err := controller.ViewMessage(id)
		if err != nil {
			return util.ErrorResponse(c, 500, err.Error())
		}

	default:
		return util.ErrorResponse(c, 400, "type is invalid")
	}
	return util.SuccessResponse(c, nil)
}

func LikePost(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	typ := c.Query("type", "article")
	switch typ {
	case "article":
		err := controller.LikeArticle(id)
		if err != nil {
			return util.ErrorResponse(c, 500, err.Error())
		}
	case "sharing":
		err := controller.LikeSharing(id)
		if err != nil {
			return util.ErrorResponse(c, 500, err.Error())
		}
	case "message":
		err := controller.LikeMessage(id)
		if err != nil {
			return util.ErrorResponse(c, 500, err.Error())
		}
	default:
		return util.ErrorResponse(c, 400, "type is invalid")
	}

	return util.SuccessResponse(c, nil)
}
