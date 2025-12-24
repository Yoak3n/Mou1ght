package entity

import (
	"strings"

	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/instance"
)

type SharingEntity struct {
	ID          string             `json:"id"`
	Content     string             `json:"content"`
	Author      UserEntity         `json:"author"`
	Tags        []PostSign         `json:"tags"`
	State       PostState          `json:"state"`
	Time        PostTimeInfo       `json:"time"`
	Attachments []AttachmentEntity `json:"attachments"`
}

func NewSharingEntityFromTable(sharing *table.SharingTable) *SharingEntity {
	user, err := instance.UseDatabase().GetUser(sharing.AuthorID)
	if err != nil {
		return nil
	}
	length := util.MeasureArticleLength(sharing.Content)
	e := &SharingEntity{
		ID:      sharing.ID,
		Content: sharing.Content,
		Author:  *NewUserEntityFromTable(user, false),
		State: PostState{
			Like:   sharing.Like,
			View:   sharing.View,
			Length: length,
		},
		Time: PostTimeInfo{
			CreatedAt: sharing.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: sharing.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		Attachments: NewAttachmentsEntityFromPaths(strings.Split(sharing.Attachment, ",")),
	}
	tags, err := instance.UseDatabase().QueryTagsByID(sharing.ID, instance.SharingTag)
	if err == nil {
		e.Tags = NewTagsInformationEntityFromTable(tags)
	}
	return e
}

func NewSharingsEntityFromTables(sharings []*table.SharingTable) []SharingEntity {
	entities := make([]SharingEntity, 0, len(sharings))
	for _, sharing := range sharings {
		entity := NewSharingEntityFromTable(sharing)
		if entity != nil {
			entities = append(entities, *entity)
		}
	}
	return entities
}
