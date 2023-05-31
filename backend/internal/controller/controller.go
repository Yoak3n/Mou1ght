package controller

import (
	"Mou1ght-Server/internal/database"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = database.GetDB()
}
