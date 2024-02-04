package dto

import "Mou1ght-Server/internal/model"

type UserDTO struct {
	Name     uint     `json:"name"`
	NickName string   `json:"nick_name"`
	Email    string   `json:"email"`
	Avatar   string   `json:"avatar"`
	Desc     string   `json:"desc"`
	Roles    []string `json:"roles"`
}

func ToUserDTO(user *model.User) UserDTO {
	return UserDTO{
		Name:     user.Name,
		Avatar:   user.Avatar,
		Desc:     user.Desc,
		NickName: user.NickName,
		Email:    user.Email,
		Roles:    user.Roles,
	}
}
