package entity

import (
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/instance"
)

type MessageEntity struct {
	ID       string          `json:"id"`
	Content  string          `json:"content"`
	Position MessagePosition `json:"position"`
	State    PostState       `json:"state"`
	Time     PostTimeInfo    `json:"time"`
	AuthorIP string          `json:"-"`
}

type MessagePosition struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

func NewMessageEntityFromTable(msg *table.MessageTable) *MessageEntity {
	length := util.MeasureArticleLength(msg.Content)
	viewDelta, likeDelta := instance.UseDatabase().GetCounterDelta("message", msg.ID)
	return &MessageEntity{
		ID:      msg.ID,
		Content: msg.Content,
		Position: MessagePosition{
			X: msg.X,
			Y: msg.Y,
			Z: msg.Z,
		},
		State: PostState{
			Like:   msg.Like + likeDelta,
			View:   msg.View + viewDelta,
			Length: length,
			Status: StatusIntToString(msg.Status),
		},
		AuthorIP: msg.AuthorIP,
		Time: PostTimeInfo{
			CreatedAt: msg.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: msg.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}
}

func NewMessagesEntityFromTables(msgs []*table.MessageTable) []*MessageEntity {
	entities := make([]*MessageEntity, 0, len(msgs))
	for _, msg := range msgs {
		entity := NewMessageEntityFromTable(msg)
		entities = append(entities, entity)
	}
	return entities
}
