package controller

import (
	"Mou1ght-Server/internal/model"
	"gorm.io/gorm"
)

func CheckExistArticle(a *model.Article, id uint) (bool, *gorm.DB) {
	result := db.First(a, id)
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, result
}
