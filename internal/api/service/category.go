package service

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/notify"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/interfaces"
	"errors"
)

type CategoryService struct {
	categories    interfaces.CategoryRepository
	categoryLinks interfaces.CategoryLinkRepository
}

func NewCategoryService(categoryRepo interfaces.CategoryRepository, categoryLinkRepo interfaces.CategoryLinkRepository) *CategoryService {
	return &CategoryService{categories: categoryRepo, categoryLinks: categoryLinkRepo}
}

func (s *CategoryService) CreateCategory(data request.CategoryRequest) error {
	if data.Label == "" {
		return errors.New("label is required")
	}
	record := &table.CategoryTable{
		ID:    util.GenCategoryID(),
		Label: data.Label,
	}
	if data.Parent != nil && *data.Parent != "" {
		if err := s.validateParent("", *data.Parent); err != nil {
			return err
		}
		record.ParentID = *data.Parent
	}
	if err := s.categories.CreateCategory(record); err != nil {
		return err
	}
	notify.RevalidateClient()
	return nil
}

func (s *CategoryService) UpdateCategory(categoryID string, data request.CategoryRequest) error {
	if categoryID == "" {
		return errors.New("category id is empty")
	}
	if data.Label == "" {
		return errors.New("label is required")
	}
	fields := map[string]any{"label": data.Label}
	if data.Parent != nil {
		parent := *data.Parent
		if parent != "" {
			if err := s.validateParent(categoryID, parent); err != nil {
				return err
			}
		}
		fields["parent_id"] = parent
	}
	if err := s.categories.UpdateCategoryFields(categoryID, fields); err != nil {
		return err
	}
	notify.RevalidateClient()
	return nil
}

// validateParent 校验目标父分类存在，且不会形成环（挂到自己或子孙下）
func (s *CategoryService) validateParent(categoryID, parentID string) error {
	if parentID == categoryID {
		return errors.New("不能将分类挂到自己下面")
	}
	all, err := s.categories.GetAllCategories()
	if err != nil {
		return err
	}
	parentMap := make(map[string]string, len(all))
	exists := false
	for _, c := range all {
		parentMap[c.ID] = c.ParentID
		if c.ID == parentID {
			exists = true
		}
	}
	if !exists {
		return errors.New("父分类不存在")
	}
	cur := parentID
	for cur != "" {
		if cur == categoryID {
			return errors.New("不能将分类挂到自己的子分类下")
		}
		next, ok := parentMap[cur]
		if !ok {
			break
		}
		cur = next
	}
	return nil
}

func (s *CategoryService) DeleteCategory(categoryID string) error {
	if categoryID == "" {
		return errors.New("category id is empty")
	}
	if err := s.categories.DeleteCategory(categoryID); err != nil {
		return err
	}
	notify.RevalidateClient()
	return nil
}

func (s *CategoryService) CategoryList() []table.CategoryTable {
	categories, err := s.categories.GetAllCategories()
	if err != nil {
		return nil
	}
	return categories
}

// CategoryListWithArticle 根据请求参数获取分类列表及其包含的文章
/**
 * 根据请求参数获取分类列表及其包含的文章
 * @param req 包含过滤条件和关键词的请求结构体
 * @return map[string]any 包含分类列表的响应结果
 */
func (s *CategoryService) CategoryListWithArticle(req *request.PostListRequest) (map[string]table.CategoryTable, []table.CategoryLinkTable) {
	// 根据关键字从数据库获取分类链接信息
	categories, links, err := s.categoryLinks.GetCategoryLinkByKeyword(req.Data.Keyword)
	if err != nil {
		return nil, nil
	}
	return categories, links
}

func (s *CategoryService) CreateCategoriesLinkToArticle(categories []string, article string) error {
	for _, category := range categories {
		lid := util.GenCategoryID()
		record := &table.CategoryLinkTable{
			ID:         lid,
			CategoryID: category,
			ArticleID:  article,
		}
		err := s.categoryLinks.CreateCategoryLink(record)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *CategoryService) GetArticlesFromCategoryLink(categoryLink *table.CategoryLinkTable, descend bool) ([]table.ArticleTable, error) {
	return s.categoryLinks.GetArticlesFromCategoryLink(categoryLink, descend)
}
