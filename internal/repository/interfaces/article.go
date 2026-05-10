package interfaces

import (
	"Mou1ght/internal/domain/model/table"
	"time"
)

type ArticleRepository interface {
	CreateArticle(article *table.ArticleTable) error
	UpdateArticle(article *table.ArticleTable) error
	AddViewCountArticle(id string) error
	AddLikeCountArticle(id string) error
	GetArticleByID(id string) (*table.ArticleTable, error)
	DeleteArticleByID(id string) error
	GetArticlesByAuthorID(authorID string, desc bool) ([]table.ArticleTable, error)
	GetArticlesByAuthorIDs(ids []string, desc bool) ([]*table.ArticleTable, error)
	GetArticles(startDate, endDate *time.Time) ([]*table.ArticleTable, error)
}
