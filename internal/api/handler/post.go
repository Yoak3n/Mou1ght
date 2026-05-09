package handler

import (
	"Mou1ght/internal/api/controller"
	"Mou1ght/internal/domain/entity"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/pkg/util"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func postTypeFromRouteName(routeName string) string {
	i := strings.Index(routeName, ".")
	if i <= 0 {
		return routeName
	}
	return routeName[:i]
}

func filterPublished(result map[string]any) {
	if result == nil {
		return
	}

	if v, ok := result["articles"]; ok {
		switch items := v.(type) {
		case []*entity.ArticleEntity:
			filtered := make([]*entity.ArticleEntity, 0, len(items))
			for _, it := range items {
				if it != nil && it.State.Status == "published" {
					filtered = append(filtered, it)
				}
			}
			result["articles"] = filtered
		}
	}

	if v, ok := result["sharings"]; ok {
		switch items := v.(type) {
		case []entity.SharingEntity:
			filtered := make([]entity.SharingEntity, 0, len(items))
			for _, it := range items {
				if it.State.Status == "published" {
					filtered = append(filtered, it)
				}
			}
			result["sharings"] = filtered
		}
	}

	if v, ok := result["messages"]; ok {
		switch items := v.(type) {
		case []*entity.MessageEntity:
			filtered := make([]*entity.MessageEntity, 0, len(items))
			for _, it := range items {
				if it != nil && it.State.Status == "published" {
					filtered = append(filtered, it)
				}
			}
			result["messages"] = filtered
		}
	}

	if v, ok := result["categories"]; ok {
		switch items := v.(type) {
		case []*entity.CategoryWithArticlesEntity:
			for _, cat := range items {
				filtered := make([]entity.ArticleEntity, 0, len(cat.Articles))
				for _, a := range cat.Articles {
					if a.State.Status == "published" {
						filtered = append(filtered, a)
					}
				}
				cat.Articles = filtered
			}
			result["categories"] = items
		}
	}

	if v, ok := result["tags"]; ok {
		switch items := v.(type) {
		case []*entity.TagWithArticlesEntity:
			for _, t := range items {
				filtered := make([]entity.ArticleEntity, 0, len(t.Articles))
				for _, a := range t.Articles {
					if a.State.Status == "published" {
						filtered = append(filtered, a)
					}
				}
				t.Articles = filtered
			}
			result["tags"] = items
		case []*entity.TagWithSharingEntity:
			for _, t := range items {
				filtered := make([]entity.SharingEntity, 0, len(t.Sharings))
				for _, s := range t.Sharings {
					if s.State.Status == "published" {
						filtered = append(filtered, s)
					}
				}
				t.Sharings = filtered
			}
			result["tags"] = items
		}
	}
}

func ListPost(c *fiber.Ctx) error {
	name := postTypeFromRouteName(c.Route().Name)
	req := &request.PostListRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	resultMap := make(map[string]any)
	// 除all外暂时未支持date_range，看需求是否需要
	switch req.Filter.Typ {
	case "category":
		resultMap = controller.CategoryListWithArticle(req)
	case "tag":
		resultMap = controller.TagListWithPost(req, name)
	case "author":
		resultMap = controller.AuthorListWithPost(req)
	case "single":
		resultMap = controller.SingleListWithPost(req, name)
	case "all":
		resultMap = controller.AllListWithPost(req)
	default:
		return util.ErrorResponse(c, 400, "Invalid filter type")
	}

	return util.SuccessResponse(c, resultMap)
}

func ListPostPublic(c *fiber.Ctx) error {
	name := postTypeFromRouteName(c.Route().Name)
	req := &request.PostListRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	resultMap := make(map[string]any)
	switch req.Filter.Typ {
	case "category":
		resultMap = controller.CategoryListWithArticle(req)
	case "tag":
		resultMap = controller.TagListWithPost(req, name)
	case "author":
		resultMap = controller.AuthorListWithPost(req)
	case "single":
		resultMap = controller.SingleListWithPost(req, name)
	case "all":
		resultMap = controller.AllListWithPost(req)
	default:
		return util.ErrorResponse(c, 400, "Invalid filter type")
	}

	filterPublished(resultMap)
	return util.SuccessResponse(c, resultMap)
}

func UpdatePostStatus(c *fiber.Ctx) error {
	req := &request.UpdatePostStatusRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	err = controller.UpdatePostStatus(req)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil)
}

func ViewPost(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	typ := c.Query("type", "article")
	log.Printf("ViewPost %s %s\n", typ, id)
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
