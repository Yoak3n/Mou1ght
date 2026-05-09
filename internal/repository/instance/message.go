package instance

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"time"

	"gorm.io/gorm"
)

func (d *Database) CreateMessage(msg *table.MessageTable) error {
	return d.DB.Create(&msg).Error
}

func (d *Database) UpdateMessage(msg *table.MessageTable) error {
	return d.DB.Save(&msg).Error
}

func (d *Database) UpdateMessagePosition(id string, pos request.MessagePosition, authorIP string, isAdmin bool) error {
	query := d.DB.Model(&table.MessageTable{}).Where("id = ?", id)
	if !isAdmin {
		query = query.Where("author_ip = ?", authorIP)
	}
	return query.Updates(map[string]interface{}{
		"x": pos.X,
		"y": pos.Y,
		"z": pos.Z,
	}).Error
}

func (d *Database) AddViewCountMessage(id string) error {
	d.BumpView("message", id, 1)
	return nil
}

func (d *Database) AddLikeCountMessage(id string) error {
	d.BumpLike("message", id, 1)
	return nil
}

func (d *Database) GetMessageByID(id string) (*table.MessageTable, error) {
	msg := &table.MessageTable{}
	err := d.DB.Where("id = ?", id).First(&msg).Error
	return msg, err
}

func (d *Database) DeleteMessageByID(id string) error {
	return d.DB.Where("id = ?", id).Delete(&table.MessageTable{}).Error
}

func (d *Database) GetMessages(startDate, endDate *time.Time) ([]*table.MessageTable, error) {
	msgs := make([]*table.MessageTable, 0)
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
	err := query.Order("created_at DESC").Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (d *Database) GetOwnedMessageIDs(authorIP string) ([]string, error) {
	var ids []string
	err := d.DB.Model(&table.MessageTable{}).Where("author_ip = ?", authorIP).Pluck("id", &ids).Error
	return ids, err
}
