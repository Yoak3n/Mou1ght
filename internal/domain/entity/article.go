package entity

import (
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/instance"
)

type ArticleEntity struct {
	ID      string       `json:"id"`
	Title   string       `json:"title"`
	Content string       `json:"content"`
	Author  UserEntity   `json:"author"`
	State   PostState    `json:"state"`
	Time    PostTimeInfo `json:"time"`
}

func NewArticleEntityFromTable(article *table.ArticleTable, detail bool) *ArticleEntity {
	user, err := instance.UseDatabase().GetUser(article.AuthorID)
	if err != nil {
		return nil
	}
	length := util.MeasureArticleLength(article.Content)
	content := ""
	if detail {
		content = article.Content
	} else {
		content = util.GenerateBriefFromMarkdown(article.Content)
	}
	return &ArticleEntity{
		ID:      article.ID,
		Title:   article.Title,
		Content: content,
		State: PostState{
			View:   article.View,
			Like:   article.Like,
			Length: length,
			Status: article.Status,
		},
		Time: PostTimeInfo{
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		},
		Author: *NewUserEntityFromTable(user, false),
	}
}
