package table

import (
	"time"

	"gorm.io/gorm"
)

type SharingAttachmentLinkTable struct {
	ID           string `gorm:"primary_key;not null;"`
	SharingID    string `gorm:"not null;index"`
	AttachmentID string `gorm:"not null;index"`
	Sort         int    `gorm:"not null;default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
