package interfaces

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"time"
)

type MessageRepository interface {
	CreateMessage(record *table.MessageTable) error
	UpdateMessage(msg *table.MessageTable) error
	UpdateMessagePosition(id string, pos request.MessagePosition, authorIP string, isAdmin bool) error
	AddViewCountMessage(id string) error
	AddLikeCountMessage(id string) error
	GetMessageByID(id string) (*table.MessageTable, error)
	DeleteMessageByID(id string) error
	GetMessages(startDate, endDate *time.Time) ([]*table.MessageTable, error)
	GetOwnedMessageIDs(authorIP string) ([]string, error)
}
