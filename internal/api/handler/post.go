package handler

import (
	"Mou1ght/internal/api/service"
	"Mou1ght/internal/domain/entity"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	articleService  *service.ArticleService
	categoryService *service.CategoryService
	sharingService  *service.SharingService
	messageService  *service.MessageService
	tagService      *service.TagService
	userService     *service.UserService
	postService     *service.PostService
	dtoService      *service.DTOService
}

func NewPostHandler(articleService *service.ArticleService, categoryService *service.CategoryService, sharingService *service.SharingService, messageService *service.MessageService, tagService *service.TagService, userService *service.UserService, postService *service.PostService, dtoService *service.DTOService) *PostHandler {
	return &PostHandler{articleService: articleService, categoryService: categoryService, sharingService: sharingService, messageService: messageService, tagService: tagService, userService: userService, postService: postService, dtoService: dtoService}
}

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

func (h *PostHandler) categories(cm map[string]table.CategoryTable, links []table.CategoryLinkTable, descend bool) map[string]any {
	resultMap := make(map[string]any)
	resultMap["categories"] = make([]*entity.CategoryWithArticlesEntity, 0)
	// 遍历分类链接，获取每个分类下的文章
	for i := range links {
		// 根据排序方式获取分类下的文章
		articles, err := h.categoryService.GetArticlesFromCategoryLink(&links[i], descend)
		if err != nil {
			continue
		}
		// 获取分类记录
		categoryRecord := cm[links[i].CategoryID]
		// 创建包含文章信息的分类实体
		category := h.dtoService.GetCategoryWithArticlesEntityFromTable(&categoryRecord, articles)
		// 将分类实体添加到返回结果中
		resultMap["categories"] = append(resultMap["categories"].([]*entity.CategoryWithArticlesEntity), category)
	}
	return resultMap
}

func (h *PostHandler) returnMap(res *service.PostResult) map[string]any {
	resultMap := make(map[string]any)
	if res.Articles != nil {
		resultMap["articles"] = h.dtoService.GetArticlesEntiesFromTable(res.Articles, false)
	}
	if res.Sharings != nil {
		resultMap["sharings"] = h.dtoService.GetSharingsEntityFromTables(res.Sharings)

	}
	if res.Messages != nil {
		resultMap["messages"] = h.dtoService.GetMessagesEntityFromTables(res.Messages)
	}
	log.Println(resultMap)
	return resultMap
}

func (h *PostHandler) tags(tm map[string]table.TagTable, links []table.TagLinkTable, descend bool, typ string) map[string]any {
	ret := make(map[string]any)
	ret["tags"] = make([]*entity.TagWithArticlesEntity, 0)

	if typ == "sharing" {
		// 如果是分享类型，设置类型为分享并初始化分享标签切片
		ret["tags"] = make([]*entity.TagWithSharingEntity, 0)
	} else {
		// 否则初始化文章标签切片
		ret["tags"] = make([]*entity.TagWithArticlesEntity, 0)
	}
	// 遍历标签链接
	for i := range links {
		if typ == "sharing" {
			// 处理分享类型标签
			sharings, err := h.tagService.GetSharingFromTagLink(&links[i], descend)
			if err != nil {
				// 发生错误时跳过当前标签
				continue
			}
			// 获取标签记录并创建分享标签实体
			tagRecord := tm[links[i].TargetID]
			tag := h.dtoService.GetTagWithSharingEntityFromTable(&tagRecord, sharings)
			// 将标签添加到结果集中
			ret["tags"] = append(ret["tags"].([]*entity.TagWithSharingEntity), tag)
		} else {
			// 处理文章类型标签
			articles, err := h.tagService.GetArticlesFromTagLink(&links[i], descend)
			if err != nil {
				// 发生错误时跳过当前标签
				continue
			}
			// 获取标签记录并创建文章标签实体
			tagRecord := tm[links[i].TagID]
			tag := h.dtoService.GetTagWithArticlesEntityFromTable(&tagRecord, articles)
			// 将标签添加到结果集中
			ret["tags"] = append(ret["tags"].([]*entity.TagWithArticlesEntity), tag)
		}
	}
	return ret
}

func (h *PostHandler) authors(users []table.UserTable, descend bool) map[string]any {
	ret := make(map[string]any)
	es := make([]*entity.UserWithPostEntity, 0)
	for _, author := range users {
		articles, _ := h.articleService.GetArticlesByAuthorID(author.ID, descend)
		sharings, _ := h.sharingService.GetSharingsByAuthorID(author.ID, descend)
		e := h.dtoService.GetUserWithPostEntityFromTable(&author, sharings, articles)
		es = append(es, e)
	}
	ret["authors"] = es
	return ret
}

func (h *PostHandler) ListPost(c *fiber.Ctx) error {
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
		cm, links := h.categoryService.CategoryListWithArticle(req)
		resultMap = h.categories(cm, links, req.Filter.Sort == "desc")
	case "tag":
		tm, links := h.tagService.TagListWithPost(req, name)
		resultMap = h.tags(tm, links, req.Filter.Sort == "desc", name)
	case "author":
		users := h.userService.AuthorListWithPost(req)
		resultMap = h.authors(users, req.Filter.Sort == "desc")
	case "single":
		m := h.postService.SingleListWithPost(req, name)
		resultMap = h.returnMap(m)
	case "all":
		m := h.postService.AllListWithPost(req)
		resultMap = h.returnMap(m)
	default:
		return util.ErrorResponse(c, 400, "Invalid filter type")
	}
	return util.SuccessResponse(c, resultMap)
}

func (h *PostHandler) ListPostPublic(c *fiber.Ctx) error {
	name := postTypeFromRouteName(c.Route().Name)
	req := &request.PostListRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	resultMap := make(map[string]any)
	switch req.Filter.Typ {
	case "category":
		cm, links := h.categoryService.CategoryListWithArticle(req)
		resultMap = h.categories(cm, links, req.Filter.Sort == "desc")
	case "tag":
		tm, links := h.tagService.TagListWithPost(req, name)
		resultMap = h.tags(tm, links, req.Filter.Sort == "desc", name)
	case "author":
		users := h.userService.AuthorListWithPost(req)
		resultMap = h.authors(users, req.Filter.Sort == "desc")
	case "single":
		m := h.postService.SingleListWithPost(req, name)
		resultMap = h.returnMap(m)
	case "all":
		m := h.postService.AllListWithPost(req)
		resultMap = h.returnMap(m)
	default:
		return util.ErrorResponse(c, 400, "Invalid filter type")
	}
	filterPublished(resultMap)
	return util.SuccessResponse(c, resultMap)
}

func (h *PostHandler) UpdatePostStatus(c *fiber.Ctx) error {
	req := &request.UpdatePostStatusRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	err = h.postService.UpdatePostStatus(req)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil)
}

func (h *PostHandler) ViewPost(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	typ := c.Query("type", "article")
	switch typ {
	case "article":
		err := h.articleService.ViewArticle(id)
		if err != nil {
			return util.ErrorResponse(c, 500, err.Error())
		}
	case "sharing":
		err := h.sharingService.ViewSharing(id)
		if err != nil {
			return util.ErrorResponse(c, 500, err.Error())
		}
	case "message":
		err := h.messageService.ViewMessage(id)
		if err != nil {
			return util.ErrorResponse(c, 500, err.Error())
		}

	default:
		return util.ErrorResponse(c, 400, "type is invalid")
	}
	return util.SuccessResponse(c, nil)
}

func (h *PostHandler) LikePost(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	typ := c.Query("type", "article")
	switch typ {
	case "article":
		err := h.articleService.LikeArticle(id)
		if err != nil {
			return util.ErrorResponse(c, 500, err.Error())
		}
	case "sharing":
		err := h.sharingService.LikeSharing(id)
		if err != nil {
			return util.ErrorResponse(c, 500, err.Error())
		}
	case "message":
		err := h.messageService.LikeMessage(id)
		if err != nil {
			return util.ErrorResponse(c, 500, err.Error())
		}
	default:
		return util.ErrorResponse(c, 400, "type is invalid")
	}

	return util.SuccessResponse(c, nil)
}
