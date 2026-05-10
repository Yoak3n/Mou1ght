package interfaces

import "Mou1ght/internal/domain/model/table"

type UserRepository interface {
	CreateUser(user *table.UserTable) error
	GetUser(id string) (*table.UserTable, error)
	GetUserByName(name string) (*table.UserTable, error)
	QueryUsers(username []string) ([]table.UserTable, error)
	UpdateUser(user *table.UserTable) error
	DeleteUser(user *table.UserTable) error
}
