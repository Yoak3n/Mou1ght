package controller

import (
	"Mou1ght-Server/internal/database"
	"Mou1ght-Server/internal/model"
)

func CheckDuplicateName(user *model.User, name string) bool {
	result := database.DB.Where("name = ?", name).First(&user)
	if result.RowsAffected == 0 {
		return false
	}
	return true
}

func RegisterUser(user *model.User) error {
	result := database.DB.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
