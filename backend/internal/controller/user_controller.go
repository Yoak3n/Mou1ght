package controller

import (
	"Mou1ght-Server/internal/model"
	"gorm.io/gorm"
)

func CheckExistName(user *model.User, name string) (bool, *gorm.DB) {
	result := db.Where("name = ?", name).First(&user)
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, result
}

func RegisterUser(user *model.User) error {
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
