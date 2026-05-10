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

func (c *CategoryRepository) UpdateCategory(category *table.CategoryTable) error {
	return c.db.Save(category).Error
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
