package dto

import (
	"Mou1ght-Server/internal/model"
)

type ArticleDTO struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	Author      uint   `json:"author"`
	AuthorName  string `json:"author_name"`
	Description string `json:"description"`
}

type ArticlePostDTO struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	AuthorID    uint   `json:"author_id"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type ArticleView struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	AuthorID    uint   `json:"author_id"`
	AuthorName  string `json:"author_name"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

func ToArticleList(as []*model.Article) []ArticleView {

	articleList := make([]ArticleView, 0)
	for _, article := range as {
		articleList = append(articleList, ArticleView{
			ID:          article.ID,
			Title:       article.Title,
			AuthorID:    article.Author,
			AuthorName:  article.AuthorName,
			Category:    article.Category,
			Description: article.Description,
		})
	}
	return articleList
}

func ToArticleDTO(a *model.Article) ArticleDTO {
	return ArticleDTO{
		ID:          a.ID,
		Title:       a.Title,
		Content:     a.Content,
		Category:    a.Category,
		Author:      a.Author,
		AuthorName:  a.AuthorName,
		Description: a.Description,
	}
}
