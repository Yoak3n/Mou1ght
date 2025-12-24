package entity

import (
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/instance"
)

type ArticleEntity struct {
	ID         string       `json:"id"`
	Title      string       `json:"title"`
	Content    string       `json:"content"`
	Categories []PostSign   `json:"categories"`
	Tags       []PostSign   `json:"tags"`
	Author     UserEntity   `json:"author"`
	State      PostState    `json:"state"`
	Time       PostTimeInfo `json:"time"`
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
	e := &ArticleEntity{
		ID:      article.ID,
		Title:   article.Title,
		Content: content,
		State: PostState{
			View:   article.View,
			Like:   article.Like,
			Length: length,
			Status: StatusIntToString(article.Status),
		},
		Categories: make([]PostSign, 0),
		Tags:       make([]PostSign, 0),
		Time: PostTimeInfo{
			CreatedAt: article.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: article.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		Author: *NewUserEntityFromTable(user, false),
	}
	tags, err := instance.UseDatabase().QueryTagsByID(article.ID, instance.ArticleTag)
	if err == nil {
		e.Tags = NewTagsInformationEntityFromTable(tags)
	}
	categories, err := instance.UseDatabase().QueryCategoriesByArticleID(article.ID)
	if err == nil {
		e.Categories = NewCategoriesInformationEntityFromTable(categories)
	}
	return e
}

func NewArticleEntityFromTableList(list []*table.ArticleTable, detail bool) []*ArticleEntity {
	es := make([]*ArticleEntity, len(list))
	for i, article := range list {
		es[i] = NewArticleEntityFromTable(article, detail)
	}
	return es
}
