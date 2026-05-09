package instance

import (
	"Mou1ght/internal/domain/model/table"
	"errors"

	"gorm.io/gorm"
)

func (d *Database) UpdatePostStatus(postType string, id string, status int8) error {
	tableName := ""
	switch postType {
	case "article":
		tableName = "article_tables"
	case "message":
		tableName = "message_tables"
	case "sharing":
		tableName = "sharing_tables"
	default:
		return errors.New("invalid post type")
	}
	if status < 0 || status > 2  {
		return errors.New("status must be 0, 1, or 2")
	}

	result := d.DB.Table(tableName).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("post not found")
	}
	return nil
}

func (d *Database) applyCounterDelta(postType, id string, viewDelta, likeDelta int64) {
	if viewDelta == 0 && likeDelta == 0 {
		return
	}

	var model any
	switch postType {
	case "article":
		model = &table.ArticleTable{}
	case "sharing":
		model = &table.SharingTable{}
	case "message":
		model = &table.MessageTable{}
	default:
		return
	}

	updates := make(map[string]any, 2)
	if viewDelta != 0 {
		updates["view"] = gorm.Expr("view + ?", viewDelta)
	}
	if likeDelta != 0 {
		updates["like"] = gorm.Expr("like + ?", likeDelta)
	}
	_ = d.DB.Model(model).Where("id = ? AND status = 1", id).Updates(updates).Error
}

func (d *Database) BumpView(postType, id string, n int64) {
	if d.counter == nil {
		return
	}
	d.counter.bumpView(postType, id, n)
}

func (d *Database) BumpLike(postType, id string, n int64) {
	if d.counter == nil {
		return
	}
	d.counter.bumpLike(postType, id, n)
}

func (d *Database) GetCounterDelta(postType, id string) (viewDelta int64, likeDelta int64) {
	if d.counter == nil {
		return 0, 0
	}
	return d.counter.get(postType, id)
}
