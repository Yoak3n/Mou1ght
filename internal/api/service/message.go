package service

import (
	"Mou1ght/internal/config"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/notify"
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
		AuthorIP: req.VisitorToken,
	}
	if config.GetConfig().Blog.Board.NeedReviewed {
		record.Status = 3
	} else {
		record.Status = 1
	}
	if err := m.messages.CreateMessage(record); err != nil {
		return err
	}
	notify.RevalidateClient()
	return nil
}

func (m *MessageService) UpdateMessage(req *request.UpdateMessageRequest) error {
	existing, err := m.messages.GetMessageByID(req.ID)
	if err != nil {
		return err
	}
	if existing == nil || existing.ID == "" {
		return fiber.NewError(404, "message not found")
	}
	if existing.AuthorIP != req.VisitorToken {
		return fiber.NewError(403, "Forbidden")
	}

	status := existing.Status
	if config.GetConfig().Blog.Board.NeedReviewed {
		status = 3
	}
	record := &table.MessageTable{
		PostBase: table.PostBase{
			ID:      req.ID,
			Content: req.Content,
			Status:  status,
		},
		X: req.Position.X,
		Y: req.Position.Y,
		Z: req.Position.Z,
	}
	if err := m.messages.UpdateMessage(record); err != nil {
		return err
	}
	notify.RevalidateClient()
	return nil
}

func (m *MessageService) UpdateMessagePosition(req *request.UpdateMessagePositionRequest, isAdmin bool) error {
	return m.messages.UpdateMessagePosition(req.ID, req.Position, req.VisitorToken, isAdmin)
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
	if err := m.messages.DeleteMessageByID(id); err != nil {
		return err
	}
	notify.RevalidateClient()
	return nil
}

func (m *MessageService) ListMessages(dateRange *request.PostFilterDate, sort string, page, pageSize int) ([]*table.MessageTable, int64, error) {
	opts := request.ListOptions{
		Page:           page,
		PageSize:       pageSize,
		AscendingOrder: sort == "asc",
	}
	if dateRange != nil {
		if dateRange.StartDate != "" {
			if s, err := time.Parse("2006-01-02 15:04:05", dateRange.StartDate); err == nil {
				opts.StartDate = &s
			}
		}
		if dateRange.EndDate != "" {
			if e, err := time.Parse("2006-01-02 15:04:05", dateRange.EndDate); err == nil {
				opts.EndDate = &e
			}
		}
	}
	return m.messages.GetMessages(opts)
}

func (m *MessageService) ListMessagesPublic(dateRange *request.PostFilterDate, sort string, page, pageSize int) ([]*table.MessageTable, int64, error) {
	opts := request.ListOptions{
		Page:           page,
		PageSize:       pageSize,
		AscendingOrder: sort == "asc",
		OnlyPublished:  true,
	}
	if dateRange != nil {
		if dateRange.StartDate != "" {
			if s, err := time.Parse("2006-01-02 15:04:05", dateRange.StartDate); err == nil {
				opts.StartDate = &s
			}
		}
		if dateRange.EndDate != "" {
			if e, err := time.Parse("2006-01-02 15:04:05", dateRange.EndDate); err == nil {
				opts.EndDate = &e
			}
		}
	}
	return m.messages.GetMessages(opts)
}

func (m *MessageService) GetOwnedMessageIDs(jti string) ([]string, error) {
	return m.messages.GetOwnedMessageIDs(jti)
}
