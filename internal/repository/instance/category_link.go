package instance

import (
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"

	"gorm.io/gorm"
)

type CategoryLinkRepository struct {
	db *gorm.DB
}

func NewCategoryLinkRepository(db *gorm.DB) *CategoryLinkRepository {
	return &CategoryLinkRepository{db: db}
}

func (c *CategoryLinkRepository) CreateCategoryLink(link *table.CategoryLinkTable) error {
	return c.db.Create(link).Error
}

func (c *CategoryLinkRepository) DeleteCategoryLink(id string) error {
	return c.db.Where("id = ?", id).Delete(&table.CategoryLinkTable{}).Error
}

func (c *CategoryLinkRepository) DeleteCategoryLinkByArticleID(articleID string) error {
	return c.db.Where("article_id = ?", articleID).Delete(&table.CategoryLinkTable{}).Error
}

func (c *CategoryLinkRepository) UpdateCategoryLinks(articleID string, categoryIDs map[string]bool) error {
	currentIDs := make([]string, 0)
	err := c.db.Where("article_id = ?", articleID).Model(&table.CategoryLinkTable{}).Pluck("category_id", &currentIDs).Error
	if err != nil {
		return err
	}
	if len(categoryIDs) == 0 {
		if len(currentIDs) == 0 {
			return nil
		}
		return c.DeleteCategoryLinkByArticleID(articleID)
	}

	currentSet := make(map[string]bool, len(currentIDs))
	var lastError error
	for _, currentID := range currentIDs {
		currentSet[currentID] = true
		if _, ok := categoryIDs[currentID]; ok {
			continue
		}
		lastError = c.db.
			Where("article_id = ? AND category_id = ?", articleID, currentID).
			Delete(&table.CategoryLinkTable{}).
			Error
	}

	for categoryID := range categoryIDs {
		if currentSet[categoryID] {
			continue
		}
		link := &table.CategoryLinkTable{
			ID:         util.GenCategoryLinkID(),
			ArticleID:  articleID,
			CategoryID: categoryID,
		}
		lastError = c.CreateCategoryLink(link)
	}

	return lastError
}

func (c *CategoryLinkRepository) GetArticlesFromCategoryLink(link *table.CategoryLinkTable, desc bool) ([]table.ArticleTable, error) {
	articles := make([]table.ArticleTable, 0)
	query := c.db.Model(&table.ArticleTable{}).Where("id = ?", link.ArticleID)
	if desc {
		query = query.Order("created_at desc")
	} else {
		query = query.Order("created_at asc")
	}
	err := query.Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (c *CategoryLinkRepository) GetCategoryLinkByKeyword(keyword []string) (map[string]table.CategoryTable, []table.CategoryLinkTable, error) {
	categories := make([]table.CategoryTable, 0)

	err := c.db.Where("label in ?", keyword).Find(&categories).Error
	if err != nil {
		return nil, nil, err
	}
	categoriesIds := make([]string, len(categories))
	categoriesMap := make(map[string]table.CategoryTable)
	for i, cat := range categories {
		categoriesIds[i] = cat.ID
		categoriesMap[cat.ID] = cat
	}
	links := make([]table.CategoryLinkTable, 0)
	err = c.db.Where("category_id IN ?", categoriesIds).Find(&links).Error
	if err != nil {
		return nil, nil, err
	}

	return categoriesMap, links, nil
}

func (c *CategoryLinkRepository) CreateCategoriesLinkToArticle(categories []string, article string) error {
	for _, category := range categories {
		lid := util.GenCategoryID()
		record := &table.CategoryLinkTable{
			ID:         lid,
			CategoryID: category,
			ArticleID:  article,
		}
		err := c.CreateCategoryLink(record)
		if err != nil {
			return err
		}
	}
	return nil
}
