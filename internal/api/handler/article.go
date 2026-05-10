package handler

import (
	"Mou1ght/internal/api/service"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
)

type ArticleHandler struct {
	articleService *service.ArticleService
	dto            *service.DTOService
}

func NewArticleHandler(articleService *service.ArticleService, dto *service.DTOService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService, dto: dto}
}

func (ah *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
	req := &request.CreateArticleRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	err = ah.articleService.CreateArticle(req)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Create article successfully")
}

func (ah *ArticleHandler) DeleteArticle(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	err := ah.articleService.DeleteArticleByID(id)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Delete article successfully")
}

func (ah *ArticleHandler) UpdateArticle(c *fiber.Ctx) error {
	req := &request.UpdateArticleRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	err = ah.articleService.UpdateArticle(req)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Update article successfully")
}

func (ah *ArticleHandler) DetailArticle(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	article, err := ah.articleService.GetArticleByID(id)
	e := ah.dto.GetArticleEntityFromTable(article, true)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, e)
}
