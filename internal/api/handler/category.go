package handler

import (
	"Mou1ght/internal/api/controller"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/pkg/util"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func CreateCategory(c *fiber.Ctx) error {
	var requestCategory request.CategoryRequest
	err := c.BodyParser(&requestCategory)
	if err != nil {
		return util.ErrorResponse(c, 400, "请求参数错误")
	}
	err = controller.CreateCategory(requestCategory)
	if err != nil {
		return util.ErrorResponse(c, 400, fmt.Sprintf("创建分类失败: %s", err.Error()))
	}
	return util.SuccessResponse(c, nil)
}

func UpdateCategory(c *fiber.Ctx) error {
	var requestCategory request.CategoryRequest
	err := c.BodyParser(&requestCategory)
	if err != nil {
		return util.ErrorResponse(c, 400, "请求参数错误")
	}
	id := c.Params("id")
	err = controller.UpdateCategory(id, requestCategory)
	if err != nil {
		return util.ErrorResponse(c, 400, fmt.Sprintf("更新分类失败: %s", err.Error()))
	}
	return util.SuccessResponse(c, nil)
}

func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	err := controller.DeleteCategory(id)
	if err != nil {
		return util.ErrorResponse(c, 400, fmt.Sprintf("删除分类失败: %s", err.Error()))
	}
	return util.SuccessResponse(c, nil)
}

func GetAllCategories(c *fiber.Ctx) error {
	categories := controller.CategoryList()
	if categories == nil {
		return util.ErrorResponse(c, 400, "获取分类失败")
	}
	return util.SuccessResponse(c, categories)
}
