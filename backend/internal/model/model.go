package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

type User struct {
	Name     string `json:"name" gorm:"unique"`
	NikeName string `json:"nike_name"`
	Password string `json:"password"`
	Roles    Roles  `json:"roles" `
	gorm.Model
}

type Paper struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  User   `json:"author" gorm:"foreignkey:Name"`
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
