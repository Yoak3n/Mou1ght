package controller

import (
	"Mou1ght-Server/internal/dto"
	"Mou1ght-Server/internal/model"
	"Mou1ght-Server/package/logger"
	"errors"

	"gorm.io/gorm"
)

func CheckExistArticle(a *model.Article, id uint) (bool, *gorm.DB) {
	result := db.First(a, id)
	if result.RowsAffected == 0 {
		return false, nil
	}
	logger.Info.Println(a.Author)
	return true, result
}

// AddArticle 创建文章
func AddArticle(a *dto.ArticlePostDTO) error {

	authors := make([]model.User, 0)
	for _, uid := range a.AuthorID {
		user := GetUserByID(uid)
		authors = append(authors, *user)
	}

	if len(authors) != 0 {
		article := &model.Article{
			Title:       a.Title,
			Description: a.Description,
			Content:     a.Content,
			Author:      authors,
			Category:    a.Category,
		}
		result := db.Create(article)
		if result.Error != nil {
			return result.Error
		}
		logger.Info.Println("Add article successfully")
		return nil
	} else {
		return errors.New("unauthorized user")
	}

}

func GetArticleById(id uint) (*model.Article, error) {
	article := &model.Article{}
	result := db.First(article, id)
	if result.RowsAffected == 0 {
		return nil, errors.New("article not found")
	}
	return article, nil
}
