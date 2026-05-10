package instance

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/repository/interfaces"
	"time"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db      *gorm.DB
	counter interfaces.PostCounter
}

func NewMessageRepository(db *gorm.DB, counter interfaces.PostCounter) *MessageRepository {
	return &MessageRepository{db: db, counter: counter}
}

func (m *MessageRepository) CreateMessage(msg *table.MessageTable) error {
	return m.db.Create(&msg).Error
}

func (m *MessageRepository) UpdateMessage(msg *table.MessageTable) error {
	return m.db.Save(&msg).Error
}

func (m *MessageRepository) UpdateMessagePosition(id string, pos request.MessagePosition, authorIP string, isAdmin bool) error {
	query := m.db.Model(&table.MessageTable{}).Where("id = ?", id)
	if !isAdmin {
		query = query.Where("author_ip = ?", authorIP)
	}
	return query.Updates(map[string]interface{}{
		"x": pos.X,
		"y": pos.Y,
		"z": pos.Z,
	}).Error
}

func (m *MessageRepository) AddViewCountMessage(id string) error {
	m.counter.BumpView("message", id, 1)
	return nil
}

func (m *MessageRepository) AddLikeCountMessage(id string) error {
	m.counter.BumpLike("message", id, 1)
	return nil
}

func (m *MessageRepository) GetMessageByID(id string) (*table.MessageTable, error) {
	msg := &table.MessageTable{}
	err := m.db.Where("id = ?", id).First(&msg).Error
	return msg, err
}

func (m *MessageRepository) DeleteMessageByID(id string) error {
	return m.db.Where("id = ?", id).Delete(&table.MessageTable{}).Error
}

func (m *MessageRepository) GetMessages(startDate, endDate *time.Time) ([]*table.MessageTable, error) {
	msgs := make([]*table.MessageTable, 0)
	var query *gorm.DB
	if startDate != nil {
		if endDate == nil {
			query = m.db.Where("created_at >= ?", startDate)
		} else {
			query = m.db.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		}
	} else {
		if endDate == nil {
			query = m.db
		} else {
			query = m.db.Where("created_at <= ?", endDate)
		}
	}
	err := query.Order("created_at DESC").Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (m *MessageRepository) GetOwnedMessageIDs(authorIP string) ([]string, error) {
	var ids []string
	err := m.db.Model(&table.MessageTable{}).Where("author_ip = ?", authorIP).Pluck("id", &ids).Error
	return ids, err
}
