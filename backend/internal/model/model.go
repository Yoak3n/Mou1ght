package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

type User struct {
	Name     uint   `json:"name" gorm:"unique"`
	NickName string `json:"nick_name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Password string `json:"password"`
	Desc     string `json:"desc"`
	Roles    Roles  `json:"roles"`
	gorm.Model
}

func (u *User) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, u)
}

func (u *User) Value() (driver.Value, error) {
	return json.Marshal(u)
}

type Roles []string

func (r *Roles) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, r)
}

func (r *Roles) Value() (driver.Value, error) {
	return json.Marshal(r)
}
