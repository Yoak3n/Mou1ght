package entity

import (
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
)

type MessageEntity struct {
	ID      string       `json:"id"`
	Content string       `json:"content"`
	Position MessagePosition `json:"position"`
	State   PostState    `json:"state"`
	Time    PostTimeInfo `json:"time"`
}

type MessagePosition struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int`json:"z"`
}


func NewMessageEntityFromTable(msg *table.MessageTable) *MessageEntity {
	length := util.MeasureArticleLength(msg.Content)
	return &MessageEntity{
		ID:      msg.ID,
		Content: msg.Content,
		Position: MessagePosition{
			X: msg.X,
			Y: msg.Y,
			Z: msg.Z,
		},
		State: PostState{
			Like:   msg.Like,
			View:   msg.View,
			Length: length,
		},
		Time: PostTimeInfo{
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
		},
	}
}