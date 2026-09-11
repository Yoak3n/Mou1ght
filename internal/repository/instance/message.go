package instance

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/repository/interfaces"

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
	return m.db.Model(&table.MessageTable{}).
		Where("id = ?", msg.ID).
		Updates(map[string]any{
			"content": msg.Content,
			"x":       msg.X,
			"y":       msg.Y,
			"z":       msg.Z,
			"status":  msg.Status,
		}).Error
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

func (m *MessageRepository) GetMessages(opts request.ListOptions) ([]*table.MessageTable, int64, error) {
	msgs := make([]*table.MessageTable, 0)
	query := m.db.Model(&table.MessageTable{})
	switch {
	case opts.StartDate != nil && opts.EndDate != nil:
		query = query.Where("created_at BETWEEN ? AND ?", opts.StartDate, opts.EndDate)
	case opts.StartDate != nil:
		query = query.Where("created_at >= ?", opts.StartDate)
	case opts.EndDate != nil:
		query = query.Where("created_at <= ?", opts.EndDate)
	}
	if opts.OnlyPublished {
		query = query.Where("status = ?", 1)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "created_at DESC"
	if opts.AscendingOrder {
		order = "created_at ASC"
	}
	query = query.Order(order)
	if opts.Paginated() {
		query = query.Offset(opts.Offset()).Limit(opts.PageSize)
	}
	if err := query.Find(&msgs).Error; err != nil {
		return nil, 0, err
	}
	return msgs, total, nil
}

func (m *MessageRepository) GetOwnedMessageIDs(authorIP string) ([]string, error) {
	var ids []string
	err := m.db.Model(&table.MessageTable{}).Where("author_ip = ?", authorIP).Pluck("id", &ids).Error
	return ids, err
}
