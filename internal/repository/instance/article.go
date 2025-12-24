package instance

import (
	"Mou1ght/internal/domain/model/table"
	"time"

	"gorm.io/gorm"
)

//func (d *Database) GetArticleByCategoryLink(links []table.CategoryLinkTable, desc bool) {
//	articles := make([]table.ArticleTable, 0)
//}

func (d *Database) CreateArticle(article *table.ArticleTable) error {
	return d.DB.Create(&article).Error
}

func (d *Database) UpdateArticle(article *table.ArticleTable) error {
	return d.DB.Save(&article).Error
}

func (d *Database) AddViewCountArticle(id string) error {
	return d.DB.Where("id = ? AND status = 1", id).Update("view", gorm.Expr("view + 1")).Error
}

func (d *Database) AddLikeCountArticle(id string) error {
	return d.DB.Where("id = ? AND status = 1", id).Update("like", gorm.Expr("like + 1")).Error
}

func (d *Database) GetArticleByID(id string) (*table.ArticleTable, error) {
	article := &table.ArticleTable{}
	err := d.DB.Where("id = ?", id).First(&article).Error
	return article, err
}

func (d *Database) DeleteArticleByID(id string) error {
	err := d.DB.Where("id = ?", id).Delete(&table.ArticleTable{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) GetArticlesByAuthorID(id string, desc bool) ([]table.ArticleTable, error) {
	articles := make([]table.ArticleTable, 0)
	order := "created_at ASC"
	if desc {
		order = "created_at DESC"
	}
	err := d.DB.Where("author_id = ?", id).Order(order).Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (d *Database) GetArticlesByAuthorIDs(ids []string, desc bool) ([]*table.ArticleTable, error) {
	articles := make([]*table.ArticleTable, 0)
	order := "created_at ASC"
	if desc {
		order = "created_at DESC"
	}
	err := d.DB.Find(&articles).Where("user_id IN ?", ids).Order(order).Error
	if err != nil {
		return nil, err
	}
	return articles, err
}

func (d *Database) GetArticles(startDate, endDate *time.Time) ([]*table.ArticleTable, error) {
	articles := make([]*table.ArticleTable, 0)
	var query *gorm.DB
	if startDate != nil {
		if endDate == nil {
			query = d.DB.Where("created_at >= ?", startDate)
		} else {
			query = d.DB.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		}
	} else {
		if endDate == nil {
			query = d.DB
		} else {
			query = d.DB.Where("created_at <= ?", endDate)
		}
	}
	err := query.Order("created_at DESC").Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}
