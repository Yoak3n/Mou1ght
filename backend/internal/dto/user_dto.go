package dto

import "Mou1ght-Server/internal/model"

type UserDTO struct {
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

func ToUserDTO(user *model.User) UserDTO {
	return UserDTO{
		Name:  user.Name,
		Email: user.Email,
		Roles: user.Roles,
	}
}
