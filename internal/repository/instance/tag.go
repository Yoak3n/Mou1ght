package instance

import (
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"

	"maps"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// CreateTag 创建标签
func (t *TagRepository) CreateTag(tag *table.TagTable) error {
	return t.db.Create(tag).Error
}

// DeleteTag 删除标签
func (t *TagRepository) DeleteTag(tagID string) error {
	return t.db.Delete(&table.TagTable{ID: tagID}).Error
}

// GetAllTags 获取所有标签
func (t *TagRepository) GetAllTags() ([]table.TagTable, error) {
	tags := make([]table.TagTable, 0)
	err := t.db.Find(&tags).Error
	return tags, err
}

// GetTagByID 根据标签ID获取标签
func (t *TagRepository) GetTagByID(id string) (*table.TagTable, error) {
	tag := &table.TagTable{}
	err := t.db.Where("id = ?", id).First(tag).Error
	return tag, err
}

// GetTagsByID 根据标签ID列表获取标签
func (t *TagRepository) GetTagsByID(ids []string) ([]table.TagTable, error) {
	tags := make([]table.TagTable, 0)
	err := t.db.Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}

// UpdateTag 更新标签
func (t *TagRepository) UpdateTag(tag *table.TagTable) error {
	return t.db.Save(tag).Error
}

// UpdateTargetLinks 更新目标的标签链接
// ids 中为 true 的标签会被添加，为 false 的标签会被删除
func (t *TagRepository) UpdateTargetLinks(targetID string, targetType table.TagType, ids map[string]bool) error {
	if len(ids) == 0 {
		return nil
	}
	currentIDs := make([]string, 0)
	unhandledIDs := make(map[string]bool)
	maps.Copy(unhandledIDs, ids)
	err := t.db.Where("target_id = ? AND target_type = ?", targetID, targetType).Model(&table.TagLinkTable{}).Pluck("tag_id", &currentIDs).Error
	if err != nil {
		return err
	}
	var lastError error
	for _, currentID := range currentIDs {
		unhandledIDs[currentID] = false
		if _, ok := ids[currentID]; !ok {
			lastError = t.DeleteTagLinkByTagID(currentID)
			if lastError != nil {
				continue
			}
		}
	}
	for k, v := range unhandledIDs {
		if !v {
			link := &table.TagLinkTable{
				TargetID:   targetID,
				TargetType: targetType,
				TagID:      k,
				ID:         util.GenTagLinkID(),
			}
			lastError = t.CreateTagLink(link)
			if lastError != nil {
				continue
			}
		}
	}
	return lastError
}

// DeleteTagWithLink 删除标签及其关联的标签链接
func (t *TagRepository) DeleteTagWithLink(tagID string) error {
	err := t.db.Where("id = ?", tagID).Delete(&table.TagTable{}).Error
	if err != nil {
		return err
	}
	return t.db.Where("tag_id = ? ", tagID).Delete(&table.TagLinkTable{}).Error
}

// CreateTagLink 创建标签链接
func (t *TagRepository) CreateTagLink(tagLink *table.TagLinkTable) error {
	return t.db.Create(tagLink).Error
}

// QueryTagsByLabel 根据标签标签列表查询标签
func (t *TagRepository) QueryTagsByLabel(label []string) ([]table.TagTable, error) {
	tags := make([]table.TagTable, 0)
	err := t.db.Where("label IN ?", label).Find(&tags).Error
	return tags, err
}

// QueryTagsByID 根据目标ID和目标类型查询标签
func (t *TagRepository) QueryTagsByID(targetID string, targetType table.TagType) ([]table.TagTable, error) {
	ids := make([]string, 0)
	err := t.db.Where("target_id = ? AND target_type = ?", targetID, targetType).Model(&table.TagLinkTable{}).Pluck("tag_id", &ids).Error
	if err != nil {
		return nil, err
	}
	return t.GetTagsByID(ids)
}

// DeleteTagLink 删除标签链接
func (t *TagRepository) DeleteTagLink(linkID string) error {
	return t.db.Where("id = ?", linkID).Delete(&table.TagLinkTable{}).Error
}

// DeleteTagLinkByTagID 删除标签关联的标签链接
func (t *TagRepository) DeleteTagLinkByTagID(tagID string) error {
	return t.db.Where("tag_id = ?", tagID).Delete(&table.TagLinkTable{}).Error
}

// DeleteTagLinkFromTarget 删除目标关联的标签链接
func (t *TagRepository) DeleteTagLinkFromTarget(targetID string, targetType table.TagType) error {
	return t.db.Where("target_id = ? AND target_type = ?", targetID, targetType).Delete(&table.TagLinkTable{}).Error
}

// GetTagLinkByKeyword 根据标签标签列表查询标签链接
func (t *TagRepository) GetTagLinkByKeyword(keyword []string, typ string) (map[string]table.TagTable, []table.TagLinkTable, error) {
	tags := make([]table.TagTable, 0)

	err := t.db.Where("label IN ?", keyword).Find(&tags).Error
	if err != nil {
		return nil, nil, err
	}
	tagsMap := make(map[string]table.TagTable)
	tagsIds := make([]string, len(tags))
	for i, tag := range tags {
		tagsMap[tag.ID] = tag
		tagsIds[i] = tag.ID
	}
	links := make([]table.TagLinkTable, 0)
	tType := 1
	if typ == "sharing" {
		tType = 2
	}
	err = t.db.Where("tag_id IN ? AND target_type = ?", tagsIds, tType).Find(&links).Error
	if err != nil {
		return nil, nil, err
	}
	return tagsMap, links, nil
}

// GetArticlesFromTagLink 根据标签链接查询文章
func (t *TagRepository) GetArticlesFromTagLink(link *table.TagLinkTable, desc bool) ([]table.ArticleTable, error) {
	articles := make([]table.ArticleTable, 0)
	query := t.db.Model(&table.ArticleTable{}).Where("id = ?", link.TargetID)
	if desc {
		query = query.Order("created_at DESC")
	} else {
		query = query.Order("created_at asc")
	}
	err := query.Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// GetSharingFromTagLink 根据标签链接查询分享
func (t *TagRepository) GetSharingFromTagLink(link *table.TagLinkTable, desc bool) ([]table.SharingTable, error) {
	sharings := make([]table.SharingTable, 0)
	query := t.db.Model(&table.SharingTable{}).Where("id = ?", link.TargetID)
	if desc {
		query = query.Order("created_at DESC")
	} else {
		query = query.Order("created_at asc")
	}
	err := query.Find(&sharings).Error
	if err != nil {
		return nil, err
	}
	return sharings, nil
}

func (t *TagRepository) CreateTagsLinkToArticle(tags []string, articleID string) error {
	for _, tag := range tags {
		lid := util.GenTagLinkID()
		record := &table.TagLinkTable{
			ID:         lid,
			TargetID:   articleID,
			TargetType: 1,
			TagID:      tag,
		}
		err := t.CreateTagLink(record)
		if err != nil {
			return err
		}
	}
	return nil
}
