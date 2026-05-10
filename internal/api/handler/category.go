package handler

import (
	"Mou1ght/internal/api/service"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/instance"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
	dto             *service.DTOService
}

func NewCategoryHandler(categoryService *service.CategoryService, categoryRepository *instance.CategoryRepository, dto *service.DTOService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService, dto: dto}
}

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var requestCategory request.CategoryRequest
	err := c.BodyParser(&requestCategory)
	if err != nil {
		return util.ErrorResponse(c, 400, "请求参数错误")
	}
	err = h.categoryService.CreateCategory(requestCategory)
	if err != nil {
		return util.ErrorResponse(c, 400, fmt.Sprintf("创建分类失败: %s", err.Error()))
	}
	return util.SuccessResponse(c, nil)
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	var requestCategory request.CategoryRequest
	err := c.BodyParser(&requestCategory)
	if err != nil {
		return util.ErrorResponse(c, 400, "请求参数错误")
	}
	id := c.Params("id")
	err = h.categoryService.UpdateCategory(id, requestCategory)
	if err != nil {
		return util.ErrorResponse(c, 400, fmt.Sprintf("更新分类失败: %s", err.Error()))
	}
	return util.SuccessResponse(c, nil)
}

func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.categoryService.DeleteCategory(id)
	if err != nil {
		return util.ErrorResponse(c, 400, fmt.Sprintf("删除分类失败: %s", err.Error()))
	}
	return util.SuccessResponse(c, nil)
}

func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	categories := h.categoryService.CategoryList()
	if categories == nil {
		return util.ErrorResponse(c, 400, "获取分类失败")
	}
	e := h.dto.GetCategoryGroupFromTables(categories)
	if e == nil {
		return util.ErrorResponse(c, 500, "获取分类失败")
	}
	return util.SuccessResponse(c, e)
}
