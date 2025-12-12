package entity

import "Mou1ght/internal/domain/model/table"

type CategoryWithArticlesEntity struct {
	Category CategoryInformationEntity `json:"category"`
	Articles []ArticleEntity           `json:"articles"`
}

type CategoryInformationEntity struct {
	ID    string `json:"id"`
	Label string `json:"label"`
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

func NewCategoryInformationEntityFromTable(category *table.CategoryTable) CategoryInformationEntity {
	return CategoryInformationEntity{
		ID:    category.ID,
		Label: category.Label,
	}
}
