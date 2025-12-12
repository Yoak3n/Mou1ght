package table

import (
	"time"

	"gorm.io/gorm"
)

type PostBase struct {
	ID      string `gorm:"primary_key;not null;"`
	Content string `gorm:"not null;"`
	Like    int64  `gorm:"not null;default:0"`
	View    int64  `gorm:"not null;default:0"`
	// 0 for draft, 1 for publish, 2 for hide
	Status    int8 `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ArticleTable struct {
	PostBase
	Title    string
	AuthorID string `gorm:"not null;"`
}

type SharingTable struct {
	PostBase
	AuthorID   string `gorm:"not null;"`
	Attachment string
}

type MessageTable struct {
	PostBase
	X        int
	Y        int
	Z        int
	AuthorIP string
}
