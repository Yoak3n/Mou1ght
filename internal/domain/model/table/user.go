package table

import (
	"time"

	"gorm.io/gorm"
)

type UserTable struct {
	ID        string `gorm:"primary_key;not null;"`
	UserName  string `gorm:"not null;"`
	Password  string `gorm:"not null;"`
	Avatar    string
	Email     string
	Phone     string
	Role      uint
	LastLogin time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
