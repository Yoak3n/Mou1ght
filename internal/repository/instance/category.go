package instance

import (
	"Mou1ght/internal/domain/model/table"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (c *CategoryRepository) CreateCategory(category *table.CategoryTable) error {
	return c.db.Create(category).Error
}

func (c *CategoryRepository) UpdateCategoryFields(id string, fields map[string]any) error {
	if id == "" || len(fields) == 0 {
		return nil
	}
	return c.db.Model(&table.CategoryTable{}).
		Where("id = ?", id).
		Updates(fields).Error
}

func (c *CategoryRepository) DeleteCategory(id string) error {
	err := c.db.Where("id = ?", id).Delete(&table.CategoryTable{}).Error
	if err != nil {
		return err
	}
	return c.db.Where("category_id = ?", id).Delete(&table.CategoryLinkTable{}).Error
}

func (c *CategoryRepository) GetAllCategories() ([]table.CategoryTable, error) {
	links := make([]table.CategoryTable, 0)
	err := c.db.Find(&links).Error
	return links, err
}

func (c *CategoryRepository) GetCategoriesByID(ids []string) ([]table.CategoryTable, error) {
	categories := make([]table.CategoryTable, 0)
	err := c.db.Where("id in ?", ids).Find(&categories).Error
	return categories, err
}

func (c *CategoryRepository) QueryCategoriesByArticleID(articleID string) ([]table.CategoryTable, error) {
	ids := make([]string, 0)
	err := c.db.Where("article_id = ?", articleID).Model(&table.CategoryLinkTable{}).Pluck("category_id", &ids).Error
	if err != nil {
		return nil, err
	}
	return c.GetCategoriesByID(ids)
}

// CountArticlesGroupByCategory 统计每个分类直接关联的已发布文章数。
// 关联文章已软删或非发布状态的链接不计入。
func (c *CategoryRepository) CountArticlesGroupByCategory() (map[string]int64, error) {
	type countRow struct {
		CategoryID string
		Count      int64
	}
	rows := make([]countRow, 0)
	err := c.db.Model(&table.CategoryLinkTable{}).
		Select("category_link_tables.category_id AS category_id, COUNT(*) AS count").
		Joins("JOIN article_tables ON article_tables.id = category_link_tables.article_id AND article_tables.status = 1 AND article_tables.deleted_at IS NULL").
		Group("category_link_tables.category_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int64, len(rows))
	for _, row := range rows {
		counts[row.CategoryID] = row.Count
	}
	return counts, nil
}
