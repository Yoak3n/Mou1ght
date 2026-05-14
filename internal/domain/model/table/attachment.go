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
	ID           string `gorm:"primary_key;not null;"`
	SharingID    string `gorm:"not null;index"`
	AttachmentID string `gorm:"not null;index"`
	Sort         int    `gorm:"not null;default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
