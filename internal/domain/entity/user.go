package entity

import "Mou1ght/internal/domain/model/table"

type UserEntity struct {
	ID                  string `json:"id"`
	UserName            string `json:"username"`
	Avatar              string `json:"avatar"`
	Role                string `json:"role"`
	*OutsideInformation `json:",omitempty"`
}

type OutsideInformation struct {
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	LastLogin string `json:"last_login"`
}

func NewUserEntityFromTable(user *table.UserTable, more bool) *UserEntity {
	e := &UserEntity{
		ID:       user.ID,
		UserName: user.UserName,
		Avatar:   user.Avatar,
	}
	rs := ""
	switch user.Role {
	case 0:
		rs = "admin"
	case 1:
		rs = "user"
	}
	e.Role = rs
	if more {
		e.OutsideInformation = &OutsideInformation{
			Email:     user.Email,
			Phone:     user.Phone,
			LastLogin: user.LastLogin.Format("2006-01-02 15:04:05"),
		}
	}
	return e
}

type UserWithPostEntity struct {
	Author   *UserEntity     `json:"author"`
	Sharings []SharingEntity `json:"sharings"`
	Articles []ArticleEntity `json:"articles"`
}


