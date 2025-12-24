package controller

import (
	"Mou1ght/internal/domain/entity"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/instance"
	"errors"
)

func CreateCategory(data request.CategoryRequest) error {
	record := &table.CategoryTable{
		ID:    util.GenCategoryID(),
		Label: data.Label,
	}
	if data.Parent != "" {
		record.ParentID = data.Parent
	}
	return instance.UseDatabase().CreateCategory(record)
}

func UpdateCategory(categoryID string, data request.CategoryRequest) error {
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
	return instance.UseDatabase().UpdateCategory(record)
}

func DeleteCategory(categoryID string) error {
	if categoryID == "" {
		return errors.New("category id is empty")
	}
	return instance.UseDatabase().DeleteCategory(categoryID)
}

func CategoryList() []*entity.CategoryGroup {
	categories, err := instance.UseDatabase().GetAllCategories()
	if err != nil {
		return nil
	}
	return entity.NewCategoryGroupFromTables(categories)
}

// CategoryListWithArticle 根据请求参数获取分类列表及其包含的文章
/**
 * 根据请求参数获取分类列表及其包含的文章
 * @param req 包含过滤条件和关键词的请求结构体
 * @return map[string]any 包含分类列表的响应结果
 */
func CategoryListWithArticle(req *request.PostListRequest) map[string]any {
	// 从请求中获取过滤条件
	filter := req.Filter
	// 判断是否为降序排列
	descend := filter.Sort == "desc"
	// 根据关键字从数据库获取分类链接信息
	categories, links, err := instance.UseDatabase().GetCategoryLinkByKeyword(req.Data.Keyword)
	if err != nil {
		return nil
	}
	// 初始化返回结果map
	ret := make(map[string]any)
	ret["categories"] = make([]*entity.CategoryWithArticlesEntity, 0)
	// 遍历分类链接，获取每个分类下的文章
	for i := range links {
		// 根据排序方式获取分类下的文章
		articles, err := instance.UseDatabase().GetArticlesFromCategoryLink(&links[i], descend)
		if err != nil {
			continue
		}
		// 获取分类记录
		categoryRecord := categories[links[i].CategoryID]
		// 创建包含文章信息的分类实体
		category := entity.NewCategoryWithArticlesEntityFromTable(&categoryRecord, articles)
		// 将分类实体添加到返回结果中
		ret["categories"] = append(ret["categories"].([]*entity.CategoryWithArticlesEntity), category)
	}
	return ret
}

func CreateCategoriesLinkToArticle(categories []string, article string) error {
	for _, category := range categories {
		lid := util.GenCategoryID()
		record := &table.CategoryLinkTable{
			ID:         lid,
			CategoryID: category,
			ArticleID:  article,
		}
		err := instance.UseDatabase().CreateCategoryLink(record)
		if err != nil {
			return err
		}
	}
	return nil
}
