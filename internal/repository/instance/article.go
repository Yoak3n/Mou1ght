package instance

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/repository/interfaces"

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
	return a.db.Model(&table.ArticleTable{}).
		Where("id = ?", article.ID).
		Updates(map[string]any{
			"title":     article.Title,
			"content":   article.Content,
			"author_id": article.AuthorID,
		}).Error
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

func (a *ArticleRepository) GetArticles(opts request.ListOptions) ([]*table.ArticleTable, int64, error) {
	articles := make([]*table.ArticleTable, 0)
	query := a.db.Model(&table.ArticleTable{})
	switch {
	case opts.StartDate != nil && opts.EndDate != nil:
		query = query.Where("created_at BETWEEN ? AND ?", opts.StartDate, opts.EndDate)
	case opts.StartDate != nil:
		query = query.Where("created_at >= ?", opts.StartDate)
	case opts.EndDate != nil:
		query = query.Where("created_at <= ?", opts.EndDate)
	}
	if opts.OnlyPublished {
		query = query.Where("status = ?", 1)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "created_at DESC"
	if opts.AscendingOrder {
		order = "created_at ASC"
	}
	query = query.Order(order)
	if opts.Paginated() {
		query = query.Offset(opts.Offset()).Limit(opts.PageSize)
	}
	if err := query.Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}
