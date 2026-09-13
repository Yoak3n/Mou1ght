package entity

import "Mou1ght/internal/domain/model/table"

type CategoryWithArticlesEntity struct {
	Category PostSign        `json:"category"`
	Articles []ArticleEntity `json:"articles"`
}

func NewCategoryInformationEntityFromTable(category *table.CategoryTable) PostSign {
	return PostSign{
		ID:    category.ID,
		Label: category.Label,
	}
}

type CategoryGroup struct {
	PostSign
	Parent   string           `json:"parent,omitempty"`
	Children []*CategoryGroup `json:"children"`
	// TotalCount 在 Count（直接关联）基础上额外聚合所有子分类的文章数
	TotalCount int64 `json:"total_count,omitempty"`
}
