package entity

import (
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/instance"
)

type SharingEntity struct {
	ID      string       `json:"id"`
	Content string       `json:"content"`
	Author  UserEntity   `json:"author"`
	State   PostState    `json:"state"`
	Time    PostTimeInfo `json:"time"`
}

func NewSharingEntityFromTable(sharing *table.SharingTable) *SharingEntity {
	user, err := instance.UseDatabase().GetUser(sharing.AuthorID)
	if err != nil {
		return nil
	}
	length := util.MeasureArticleLength(sharing.Content)
	return &SharingEntity{
		ID:      sharing.ID,
		Content: sharing.Content,
		Author:  *NewUserEntityFromTable(user, false),
		State: PostState{
			Like:   sharing.Like,
			View:   sharing.View,
			Length: length,
		},
		Time: PostTimeInfo{
			CreatedAt: sharing.CreatedAt,
			UpdatedAt: sharing.UpdatedAt,
		},
	}
}
