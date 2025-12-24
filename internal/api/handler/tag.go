package handler

import (
	"Mou1ght/internal/api/controller"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/pkg/util"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func CreateTag(c *fiber.Ctx) error {
	var requestTag request.CreateTagRequest
	err := c.BodyParser(&requestTag)
	if err != nil {
		return util.ErrorResponse(c, 400, "请求参数错误")
	}
	err = controller.CreateTag(&requestTag)
	if err != nil {
		return util.ErrorResponse(c, 400, fmt.Sprintf("创建标签失败: %s", err.Error()))
	}
	return util.SuccessResponse(c, nil)
}

func DeleteTag(c *fiber.Ctx) error {
	id := c.Params("id")
	err := controller.DeleteTag(id)
	if err != nil {
		return util.ErrorResponse(c, 400, fmt.Sprintf("删除标签失败: %s", err.Error()))
	}
	return util.SuccessResponse(c, nil)
}

func GetAllTags(c *fiber.Ctx) error {
	tags := controller.TagsList()
	if tags == nil {
		return util.ErrorResponse(c, 400, "获取标签失败")
	}
	return util.SuccessResponse(c, tags)
}
