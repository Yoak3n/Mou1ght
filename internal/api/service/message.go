package service

import (
	"Mou1ght/internal/config"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/interfaces"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MessageService struct {
	messages interfaces.MessageRepository
}

func NewMessageService(messages interfaces.MessageRepository) *MessageService {
	return &MessageService{messages: messages}
}

func (m *MessageService) CreateMessage(req *request.CreateMessageRequest) error {
	question := strings.TrimSpace(config.GetConfig().Blog.Board.Question)
	if question != "" {
		expected := strings.TrimSpace(config.GetConfig().Blog.Board.Answer)
		if expected == "" {
			return fiber.NewError(500, "Board question enabled but answer is not configured")
		}
		if strings.TrimSpace(req.BoardAnswer) != expected {
			return fiber.NewError(403, "Incorrect answer")
		}
	}

	mid := util.GenMessageID()
	record := &table.MessageTable{
		PostBase: table.PostBase{
			ID:      mid,
			Content: req.Content,
		},
		X:        req.Position.X,
		Y:        req.Position.Y,
		Z:        req.Position.Z,
		AuthorIP: req.AuthorIP,
	}
	if config.GetConfig().Blog.Board.NeedReviewed {
		record.Status = 3
	} else {
		record.Status = 1
	}
	return m.messages.CreateMessage(record)
}

func (m *MessageService) UpdateMessage(req *request.UpdateMessageRequest) error {
	record := &table.MessageTable{
		PostBase: table.PostBase{
			ID:      req.ID,
			Content: req.Content,
		},
		X:        req.Position.X,
		Y:        req.Position.Y,
		Z:        req.Position.Z,
		AuthorIP: req.AuthorIP,
	}
	if config.GetConfig().Blog.Board.NeedReviewed {
		record.Status = 3
	} else {
		record.Status = 1
	}
	return m.messages.UpdateMessage(record)
}

func (m *MessageService) UpdateMessagePosition(req *request.UpdateMessagePositionRequest, isAdmin bool) error {
	return m.messages.UpdateMessagePosition(req.ID, req.Position, req.AuthorIP, isAdmin)
}

func (m *MessageService) ViewMessage(id string) error {
	return m.messages.AddViewCountMessage(id)
}

func (m *MessageService) LikeMessage(id string) error {
	return m.messages.AddLikeCountMessage(id)
}

func (m *MessageService) GetMessageByID(id string) (*table.MessageTable, error) {
	record, err := m.messages.GetMessageByID(id)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (m *MessageService) DeleteMessageByID(id string) error {
	return m.messages.DeleteMessageByID(id)
}

func (m *MessageService) ListMessages(dateRange *request.PostFilterDate, sort string) ([]*table.MessageTable, error) {
	var startDate, endDate *time.Time
	if dateRange != nil {
		s, _ := time.Parse("2006-01-02 15:04:05", dateRange.StartDate)
		e, _ := time.Parse("2006-01-02 15:04:05", dateRange.EndDate)
		startDate = &s
		endDate = &e
	}
	msgs, err := m.messages.GetMessages(startDate, endDate)
	if err != nil {
		return nil, err
	}
	if sort == "desc" {
		for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
			msgs[i], msgs[j] = msgs[j], msgs[i]
		}
	}
	return msgs, nil
}

func (m *MessageService) GetOwnedMessageIDs(jti string) ([]string, error) {
	return m.messages.GetOwnedMessageIDs(jti)
}
