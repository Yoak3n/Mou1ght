package dto

import "Mou1ght-Server/internal/model"

type ArticleDTO struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Category    string   `json:"category"`
	Authors     []string `json:"author"`
	Description string   `json:"description"`
}

type ArticlePostDTO struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	AuthorID    []uint `json:"author_id"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

func ToArticleDTO(a *model.Article) ArticleDTO {

	authors := make([]string, 0)
	for _, author := range a.Author {
		authors = append(authors, author.NickName)
	}
	return ArticleDTO{
		Title:       a.Title,
		Content:     a.Content,
		Category:    a.Category,
		Authors:     authors,
		Description: a.Description,
	}
}
