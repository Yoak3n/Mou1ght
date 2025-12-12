package entity

import "Mou1ght/internal/domain/model/table"

type UserEntity struct {
	ID                  string `json:"id"`
	UserName            string `json:"username"`
	Avatar              string `json:"avatar"`
	Role                string `json:"role"`
	*OutsideInformation `json:"omitempty"`
}

type OutsideInformation struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
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
			Email: user.Email,
			Phone: user.Phone,
		}
	}
	return e
}

type UserWithPostEntity struct {
	Author   *UserEntity     `json:"author"`
	Sharings []SharingEntity `json:"sharings"`
	Articles []ArticleEntity `json:"articles"`
}

func NewUserWithPostEntityFromTable(author *table.UserTable, sharings []table.SharingTable, articles []table.ArticleTable) *UserWithPostEntity {
	sharing := make([]SharingEntity, 0)
	for _, s := range sharings {
		sharing = append(sharing, *NewSharingEntityFromTable(&s))
	}
	article := make([]ArticleEntity, 0)
	for _, a := range articles {
		article = append(article, *NewArticleEntityFromTable(&a, false))
	}

	return &UserWithPostEntity{
		Author:   NewUserEntityFromTable(author, false),
		Sharings: sharing,
		Articles: article,
	}
}
