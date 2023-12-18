package model

import "gorm.io/gorm"

type Article struct {
	Title       string `json:"title" gorm:"unique"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	Author      User   `json:"author" `
	Description string `json:"description"`
	gorm.Model
}
