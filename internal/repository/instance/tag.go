package instance

import (
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"maps"

	"gorm.io/gorm"
)

const (
	ArticleTag TagType = 1
	SharingTag TagType = 2
)

type TagType = int

func (d *Database) CreateTag(tag *table.TagTable) error {
	return d.DB.Create(tag).Error
}

func (d *Database) GetAllTags() ([]table.TagTable, error) {
	tags := make([]table.TagTable, 0)
	err := d.DB.Find(&tags).Error
	return tags, err
}

func (d *Database) GetTagByID(id string) (*table.TagTable, error) {
	tag := &table.TagTable{}
	err := d.DB.Where("id = ?", id).First(tag).Error
	return tag, err
}

func (d *Database) GetTagsByID(ids []string) ([]table.TagTable, error) {
	tags := make([]table.TagTable, 0)
	err := d.DB.Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}

func (d *Database) UpdateTag(tag *table.TagTable) error {
	return d.DB.Save(tag).Error
}

func (d *Database) UpdateTargetLinks(targetID string, targetType TagType, ids map[string]bool) error {
	if len(ids) == 0 {
		return nil
	}
	currentIDs := make([]string, 0)
	unhandledIDs := make(map[string]bool)
	maps.Copy(unhandledIDs, ids)
	err := d.DB.Where("target_id = ? AND target_type = ?", targetID, targetType).Model(&table.TagLinkTable{}).Pluck("tag_id", &currentIDs).Error
	if err != nil {
		return err
	}
	var lastError error
	for _, currentID := range currentIDs {
		unhandledIDs[currentID] = false
		if _, ok := ids[currentID]; !ok {
			lastError = d.DeleteTagLinkByTagID(currentID)
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
			lastError = d.CreateTagLink(link)
			if lastError != nil {
				continue
			}
		}
	}
	return lastError
}

func (d *Database) DeleteTagWithLink(tagID string) error {
	err := d.DB.Where("id = ?", tagID).Delete(&table.TagTable{}).Error
	if err != nil {
		return err
	}
	return d.DB.Where("tag_id = ? ", tagID).Delete(&table.TagLinkTable{}).Error
}

func (d *Database) CreateTagLink(tagLink *table.TagLinkTable) error {
	return d.DB.Create(tagLink).Error
}

func (d *Database) QueryTagsByLabel(label []string) ([]table.TagTable, error) {
	tags := make([]table.TagTable, 0)
	err := d.DB.Where("label IN ?", label).Find(&tags).Error
	return tags, err
}

func (d *Database) QueryTagsByID(targetID string, targetType TagType) ([]table.TagTable, error) {
	ids := make([]string, 0)
	err := d.DB.Where("target_id = ? AND target_type = ?", targetID, targetType).Model(&table.TagLinkTable{}).Pluck("tag_id", &ids).Error
	if err != nil {
		return nil, err
	}
	return d.GetTagsByID(ids)
}

func (d *Database) DeleteTagLink(linkID string) error {
	return d.DB.Where("id = ?", linkID).Delete(&table.TagLinkTable{}).Error
}

func (d *Database) DeleteTagLinkByTagID(tagID string) error {
	return d.DB.Where("tag_id = ?", tagID).Delete(&table.TagLinkTable{}).Error
}

func (d *Database) DeleteTagLinkFromTarget(targetID string, targetType TagType) error {
	return d.DB.Where("target_id = ? AND target_type = ?", targetID, targetType).Delete(&table.TagLinkTable{}).Error
}

func (d *Database) GetTagLinkByKeyword(keyword []string, typ string) (map[string]table.TagTable, []table.TagLinkTable, error) {
	tags := make([]table.TagTable, 0)

	err := d.DB.Where("label IN ?", keyword).Find(&tags).Error
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
	t := 1
	if typ == "sharing" {
		t = 2
	}
	err = d.DB.Where("target_id IN ? AND target_type = ?", tagsIds, t).Find(&links).Error
	if err != nil {
		return nil, nil, err
	}
	return tagsMap, links, nil
}

func (d *Database) GetArticlesFromTagLink(link *table.TagLinkTable, desc bool) ([]table.ArticleTable, error) {
	articles := make([]table.ArticleTable, 0)
	query := d.DB.Model(&table.ArticleTable{}).Preload("created_at", func(tx *gorm.DB) *gorm.DB {
		if desc {
			return tx.Order("created_at DESC")
		}
		return tx.Order("created_at asc")
	}).Where("id = ?", link.TargetID)
	err := query.Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (d *Database) GetSharingFromTagLink(link *table.TagLinkTable, desc bool) ([]table.SharingTable, error) {
	sharings := make([]table.SharingTable, 0)
	query := d.DB.Model(&table.SharingTable{}).Preload("created_at", func(tx *gorm.DB) *gorm.DB {
		if desc {
			return tx.Order("created_at DESC")
		}
		return tx.Order("created_at asc")
	}).Where("id = ?", link.TargetID)
	err := query.Find(&sharings).Error
	if err != nil {
		return nil, err
	}
	return sharings, nil
}
