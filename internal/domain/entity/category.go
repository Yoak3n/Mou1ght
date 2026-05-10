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
	Parent   string           `json:"-"`
	Children []*CategoryGroup `json:"children"`
}
