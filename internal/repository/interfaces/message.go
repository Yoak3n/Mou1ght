package interfaces

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
)

type MessageRepository interface {
	CreateMessage(record *table.MessageTable) error
	UpdateMessage(msg *table.MessageTable) error
	UpdateMessagePosition(id string, pos request.MessagePosition, authorIP string, isAdmin bool) error
	AddViewCountMessage(id string) error
	AddLikeCountMessage(id string) error
	GetMessageByID(id string) (*table.MessageTable, error)
	DeleteMessageByID(id string) error
	GetMessages(opts request.ListOptions) ([]*table.MessageTable, int64, error)
	GetOwnedMessageIDs(authorIP string) ([]string, error)
}
