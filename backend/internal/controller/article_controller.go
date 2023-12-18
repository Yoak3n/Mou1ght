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
func AddArticle(a *dto.ArticleDTO) error {
	// 创建文章
	user := GetUserById(a.Author)
	if user != nil {
		article := &model.Article{
			Title: a.Title,
			Author: model.User{
				Name: user.Name,
			},
			Description: a.Description,
			Content:     a.Content,
		}
		db.Create(article)
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
