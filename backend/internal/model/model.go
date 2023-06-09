package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

type User struct {
	Name     string ` gorm:"unique"`
	NickName string
	Email    string
	Password string
	Roles    Roles
	gorm.Model
}

type Article struct {
	Title    string `json:"title" gorm:"unique"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Author   User   `json:"author" gorm:"foreignkey:Name"`
	gorm.Model
}

type Roles []string

func (r *Roles) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, r)
}

func (r *Roles) Value() (driver.Value, error) {
	return json.Marshal(r)
}
