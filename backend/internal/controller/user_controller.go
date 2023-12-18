package controller

import (
	"Mou1ght-Server/internal/model"

	"gorm.io/gorm"
)

func GetUserByID(id uint) (user *model.User) {
	result := db.First(&user, id)
	if result.RowsAffected == 0 {
		return nil
	}
	return
}

func CheckExistName(user *model.User, name string) (*gorm.DB, bool) {
	result := db.Where("name = ?", name).First(&user)
	if result.RowsAffected == 0 {
		return nil, false
	}
	return result, true
}

func RegisterUser(user *model.User) error {

	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
