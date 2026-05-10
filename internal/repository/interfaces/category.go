package interfaces

import "Mou1ght/internal/domain/model/table"

type CategoryRepository interface {
	CreateCategory(category *table.CategoryTable) error
	UpdateCategory(category *table.CategoryTable) error
	DeleteCategory(categoryID string) error
	GetAllCategories() ([]table.CategoryTable, error)
	GetCategoriesByID(ids []string) ([]table.CategoryTable, error)
	QueryCategoriesByArticleID(articleID string) ([]table.CategoryTable, error)
}

type CategoryLinkRepository interface {
	CreateCategoryLink(link *table.CategoryLinkTable) error
	CreateCategoriesLinkToArticle(categories []string, article string) error
	DeleteCategoryLink(linkID string) error
	DeleteCategoryLinkByArticleID(articleID string) error
	UpdateCategoryLinks(articleID string, categoryIDs map[string]bool) error
	GetArticlesFromCategoryLink(link *table.CategoryLinkTable, desc bool) ([]table.ArticleTable, error)
	GetCategoryLinkByKeyword(keyword []string) (map[string]table.CategoryTable, []table.CategoryLinkTable, error)
}
