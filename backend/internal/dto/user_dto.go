package dto

import "Mou1ght-Server/internal/model"

type UserDTO struct {
	Name     uint     `json:"name"`
	NickName string   `json:"nick_name"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

func ToUserDTO(user *model.User) UserDTO {
	return UserDTO{
		NickName: user.NickName,
		Email:    user.Email,
		Roles:    user.Roles,
	}
}
