package model

import (
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title       string `json:"title" gorm:"unique"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	Author      []User `json:"author" gorm:"ForeignKey:Name" `
	Description string `json:"description"`
}
