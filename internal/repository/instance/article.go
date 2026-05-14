package instance

import (
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/repository/interfaces"
	"time"

	"gorm.io/gorm"
)

//	func (d *Database) GetArticleByCategoryLink(links []table.CategoryLinkTable, desc bool) {
//		articles := make([]table.ArticleTable, 0)
//	}
type ArticleRepository struct {
	db      *gorm.DB
	counter interfaces.PostCounter
}

func NewArticleRepository(db *gorm.DB, counter interfaces.PostCounter) *ArticleRepository {
	return &ArticleRepository{db: db, counter: counter}
}

func (a *ArticleRepository) CreateArticle(article *table.ArticleTable) error {
	return a.db.Create(article).Error
}

func (a *ArticleRepository) UpdateArticle(article *table.ArticleTable) error {
	return a.db.Save(article).Error
}

func (a *ArticleRepository) AddViewCountArticle(id string) error {
	a.counter.BumpView("article", id, 1)
	return nil
}

func (a *ArticleRepository) AddLikeCountArticle(id string) error {
	a.counter.BumpLike("article", id, 1)
	return nil
}

func (a *ArticleRepository) GetArticleByID(id string) (*table.ArticleTable, error) {
	article := &table.ArticleTable{}
	err := a.db.Where("id = ?", id).First(&article).Error
	return article, err
}

func (a *ArticleRepository) DeleteArticleByID(id string) error {
	err := a.db.Where("id = ?", id).Delete(&table.ArticleTable{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *ArticleRepository) GetArticlesByAuthorID(id string, desc bool) ([]table.ArticleTable, error) {
	articles := make([]table.ArticleTable, 0)
	order := "created_at ASC"
	if desc {
		order = "created_at DESC"
	}
	err := a.db.Where("author_id = ?", id).Order(order).Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (a *ArticleRepository) GetArticlesByAuthorIDs(ids []string, desc bool) ([]*table.ArticleTable, error) {
	articles := make([]*table.ArticleTable, 0)
	order := "created_at ASC"
	if desc {
		order = "created_at DESC"
	}
	err := a.db.Where("author_id IN ?", ids).Order(order).Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (a *ArticleRepository) GetArticles(startDate, endDate *time.Time) ([]*table.ArticleTable, error) {
	articles := make([]*table.ArticleTable, 0)
	var query *gorm.DB
	if startDate != nil {
		if endDate == nil {
			query = a.db.Where("created_at >= ?", startDate)
		} else {
			query = a.db.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		}
	} else {
		if endDate == nil {
			query = a.db
		} else {
			query = a.db.Where("created_at <= ?", endDate)
		}
	}
	err := query.Order("created_at DESC").Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}
