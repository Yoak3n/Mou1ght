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

// GetCategoryLinkByKeyword 按分类 label 查询，并展开子孙分类：
// 返回的 map 只含命中的分类本身；links 中子孙分类的 CategoryID 会归并到命中的祖先，
// 便于上层按「请求的分类」聚合出含子分类的文章列表。
func (c *CategoryLinkRepository) GetCategoryLinkByKeyword(keyword []string) (map[string]table.CategoryTable, []table.CategoryLinkTable, error) {
	categoriesMap := make(map[string]table.CategoryTable)
	if len(keyword) == 0 {
		return categoriesMap, []table.CategoryLinkTable{}, nil
	}

	all := make([]table.CategoryTable, 0)
	if err := c.db.Find(&all).Error; err != nil {
		return nil, nil, err
	}

	byLabel := make(map[string][]table.CategoryTable, len(all))
	childrenOf := make(map[string][]string, len(all))
	for _, cat := range all {
		byLabel[cat.Label] = append(byLabel[cat.Label], cat)
		if cat.ParentID != "" {
			childrenOf[cat.ParentID] = append(childrenOf[cat.ParentID], cat.ID)
		}
	}

	rootIDs := make([]string, 0)
	seenRoot := make(map[string]bool)
	for _, label := range keyword {
		for _, cat := range byLabel[label] {
			if seenRoot[cat.ID] {
				continue
			}
			seenRoot[cat.ID] = true
			rootIDs = append(rootIDs, cat.ID)
			categoriesMap[cat.ID] = cat
		}
	}
	if len(rootIDs) == 0 {
		return categoriesMap, []table.CategoryLinkTable{}, nil
	}

	// descendant/root ID -> 命中的根分类 ID
	ownerRoot := make(map[string]string)
	for _, rootID := range rootIDs {
		if _, ok := ownerRoot[rootID]; !ok {
			ownerRoot[rootID] = rootID
		}
		queue := []string{rootID}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for _, childID := range childrenOf[cur] {
				if _, visited := ownerRoot[childID]; visited {
					continue
				}
				ownerRoot[childID] = rootID
				queue = append(queue, childID)
			}
		}
	}

	scopeIDs := make([]string, 0, len(ownerRoot))
	for id := range ownerRoot {
		scopeIDs = append(scopeIDs, id)
	}

	rawLinks := make([]table.CategoryLinkTable, 0)
	if err := c.db.Where("category_id IN ?", scopeIDs).Find(&rawLinks).Error; err != nil {
		return nil, nil, err
	}

	links := make([]table.CategoryLinkTable, 0, len(rawLinks))
	for _, link := range rawLinks {
		root, ok := ownerRoot[link.CategoryID]
		if !ok {
			continue
		}
		rolled := link
		rolled.CategoryID = root
		links = append(links, rolled)
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
