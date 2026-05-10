package table

import (
	"time"

	"gorm.io/gorm"
)

const (
	ArticleTag TagType = 1
	SharingTag TagType = 2
)

type TagType = int

type TagTable struct {
	ID        string `gorm:"primary_key;not null;"`
	Label     string `gorm:"unique;not null;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type TagLinkTable struct {
	ID       string `gorm:"primary_key;not null;"`
	TagID    string `gorm:"not null;"`
	TargetID string
	// 1 for article, 2 for sharing
	TargetType int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
