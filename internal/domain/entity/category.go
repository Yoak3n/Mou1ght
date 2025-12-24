package entity

import "Mou1ght/internal/domain/model/table"

type CategoryWithArticlesEntity struct {
	Category PostSign        `json:"category"`
	Articles []ArticleEntity `json:"articles"`
}

func NewCategoryWithArticlesEntityFromTable(category *table.CategoryTable, articles []table.ArticleTable) *CategoryWithArticlesEntity {
	s := make([]ArticleEntity, len(articles))
	for i, article := range articles {
		s[i] = *NewArticleEntityFromTable(&article, false)
	}
	return &CategoryWithArticlesEntity{
		Category: NewCategoryInformationEntityFromTable(category),
		Articles: s,
	}
}

func NewCategoriesInformationEntityFromTable(categories []table.CategoryTable) []PostSign {
	s := make([]PostSign, len(categories))
	for i, category := range categories {
		s[i] = NewCategoryInformationEntityFromTable(&category)
	}
	return s
}

func NewCategoryInformationEntityFromTable(category *table.CategoryTable) PostSign {
	return PostSign{
		ID:    category.ID,
		Label: category.Label,
	}
}

type CategoryGroup struct {
	PostSign
	Parent   string           `json:"-"`
	Children []*CategoryGroup `json:"children"`
}

func NewCategoryGroupFromTables(items []table.CategoryTable) []*CategoryGroup {
	nodeMap := make(map[string]*CategoryGroup)
	for _, item := range items {
		node := &CategoryGroup{
			PostSign: NewCategoryInformationEntityFromTable(&item),
			Parent:   item.ParentID,
			Children: make([]*CategoryGroup, 0),
		}
		nodeMap[item.ID] = node
	}
	var rootNodes = make([]*CategoryGroup, 0)
	for _, item := range items {
		node := nodeMap[item.ID]
		if item.ParentID != "" {
			if parent, ok := nodeMap[item.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		} else {
			rootNodes = append(rootNodes, node)
		}
	}
	return rootNodes
}
