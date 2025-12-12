package instance

import (
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"maps"

	"gorm.io/gorm"
)

func (d *Database) CreateCategoryLink(link *table.CategoryLinkTable) error {
	return d.DB.Create(&link).Error
}

func (d *Database) DeleteCategoryLink(id string) error {
	return d.DB.Where("id = ?", id).Delete(&table.CategoryLinkTable{}).Error
}

func (d *Database) DeleteCategoryLinkByArticleID(articleID string) error {
	return d.DB.Where("article_id = ?", articleID).Delete(&table.CategoryLinkTable{}).Error
}

func (d *Database) GetAllCategories() ([]table.CategoryLinkTable, error) {
	links := make([]table.CategoryLinkTable, 0)
	err := d.DB.Find(&links).Error
	return links, err
}

func (d *Database) UpdateCategoryLinks(articleID string, categoryIDs map[string]bool) error {
	if len(categoryIDs) == 0 {
		return nil
	}
	unhandledIDs := make(map[string]bool)
	maps.Copy(unhandledIDs, categoryIDs)
	currentIDs := make([]string, 0)
	err := d.DB.Where("article_id = ?", articleID).Pluck("id", &currentIDs).Error
	if err != nil {
		return err
	}
	var lastError error
	for _, currentID := range currentIDs {
		unhandledIDs[currentID] = false
		if _, ok := categoryIDs[currentID]; !ok {
			lastError = d.DeleteCategoryLinkByArticleID(currentID)
		}
	}
	for k, v := range unhandledIDs {
		if !v {
			link := &table.CategoryLinkTable{
				ID:         util.GenCategoryID(),
				ArticleID:  articleID,
				CategoryID: k,
			}
			lastError = d.CreateCategoryLink(link)
		}
	}

	return lastError
}

func (d *Database) GetCategoryLinkByKeyword(keyword []string) (map[string]table.CategoryTable, []table.CategoryLinkTable, error) {
	categories := make([]table.CategoryTable, 0)
	categoriesIds := make([]string, 0)
	err := d.DB.Where("label in ?", keyword).Find(&categories).Error
	if err != nil {
		return nil, nil, err
	}
	categoriesMap := make(map[string]table.CategoryTable)
	for i, cat := range categories {
		categoriesIds[i] = cat.ID
		categoriesMap[cat.ID] = cat
	}
	links := make([]table.CategoryLinkTable, 0)
	err = d.DB.Where("category_id IN ?", categoriesIds).Find(&links).Error
	if err != nil {
		return nil, nil, err
	}

	return categoriesMap, links, nil
}

func (d *Database) GetArticlesFromCategoryLink(link *table.CategoryLinkTable, desc bool) ([]table.ArticleTable, error) {
	articles := make([]table.ArticleTable, 0)
	query := d.DB.Model(&table.ArticleTable{}).Preload("created_at", func(tx *gorm.DB) *gorm.DB {
		if desc {
			return tx.Order("created_at desc")
		}
		return tx.Order("created_at asc")
	}).Where("id = ?", link.ArticleID)
	err := query.Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}
