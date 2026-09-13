package table

import (
	"time"

	"gorm.io/gorm"
)

type AttachmentTable struct {
	ID           string `gorm:"primary_key;not null;"`
	OriginalName string `gorm:"not null;"`
	StoragePath  string `gorm:"not null;uniqueIndex"`
	Mime         string `gorm:"not null;"`
	Sha256       string `gorm:"not null;index"`
	Size         int64  `gorm:"not null;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type AttachmentLinkTable struct {
	ID string `gorm:"primary_key;not null;"`
	// SharingID 与 ArticleID 二选一：说说链接填 SharingID，文章链接填 ArticleID，另一侧留空
	SharingID    string `gorm:"not null;index"`
	ArticleID    string `gorm:"not null;index"`
	AttachmentID string `gorm:"not null;index"`
	Sort         int    `gorm:"not null;default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
