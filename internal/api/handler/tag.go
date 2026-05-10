package handler

import (
	"Mou1ght/internal/api/service"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/pkg/util"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type TagHandler struct {
	tagServ *service.TagService
}

func NewTagHandler(ts *service.TagService) *TagHandler {
	return &TagHandler{tagServ: ts}
}

func (th *TagHandler) CreateTag(c *fiber.Ctx) error {
	var requestTag request.CreateTagRequest
	err := c.BodyParser(&requestTag)
	if err != nil {
		return util.ErrorResponse(c, 400, "请求参数错误")
	}
	err = th.tagServ.CreateTag(&requestTag)
	if err != nil {
		return util.ErrorResponse(c, 400, fmt.Sprintf("创建标签失败: %s", err.Error()))
	}
	return util.SuccessResponse(c, nil)
}

func (th *TagHandler) DeleteTag(c *fiber.Ctx) error {
	id := c.Params("id")
	err := th.tagServ.DeleteTag(id)
	if err != nil {
		return util.ErrorResponse(c, 400, fmt.Sprintf("删除标签失败: %s", err.Error()))
	}
	return util.SuccessResponse(c, nil)
}

func (th *TagHandler) GetAllTags(c *fiber.Ctx) error {
	tags := th.tagServ.TagsList()
	if tags == nil {
		return util.ErrorResponse(c, 400, "获取标签失败")
	}
	return util.SuccessResponse(c, tags)
}
