package service

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
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
	record := &table.CategoryTable{
		ID:    util.GenCategoryID(),
		Label: data.Label,
	}
	if data.Parent != "" {
		record.ParentID = data.Parent
	}
	return s.categories.CreateCategory(record)
}

func (s *CategoryService) UpdateCategory(categoryID string, data request.CategoryRequest) error {
	if categoryID == "" {
		return errors.New("category id is empty")
	}
	record := &table.CategoryTable{
		ID:    categoryID,
		Label: data.Label,
	}
	if data.Parent != "" {
		record.ParentID = data.Parent
	}
	return s.categories.UpdateCategory(record)
}

func (s *CategoryService) DeleteCategory(categoryID string) error {
	if categoryID == "" {
		return errors.New("category id is empty")
	}
	return s.categories.DeleteCategory(categoryID)
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
