package dto

import "Mou1ght-Server/internal/model"

type ArticleDTO struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	Author      uint   `json:"author"`
	Description string `json:"description"`
}

func ToArticleDTO(a *model.Article) ArticleDTO {
	return ArticleDTO{
		Title:       a.Title,
		Content:     a.Content,
		Category:    a.Category,
		Author:      a.Author.ID,
		Description: a.Description,
	}
}
