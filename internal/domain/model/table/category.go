package table

import (
	"time"

	"gorm.io/gorm"
)

type CategoryTable struct {
	ID        string `gorm:"primary_key;not null;"`
	Label     string `gorm:"not null;"`
	ParentID  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type CategoryLinkTable struct {
	ID         string `gorm:"primary_key;not null;"`
	ArticleID  string `gorm:"not null"`
	CategoryID string `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
