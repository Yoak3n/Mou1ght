package interfaces

import "Mou1ght/internal/domain/model/table"

type UserRepository interface {
	CreateUser(user *table.UserTable) error
	GetUser(id string) (*table.UserTable, error)
	GetUserByName(name string) (*table.UserTable, error)
	QueryUsers(username []string) ([]table.UserTable, error)
	CountUsers() (int64, error)
	UpdateUser(user *table.UserTable) error
	UpdateUserProfile(id string, fields map[string]any) error
	UpdateUserPassword(id, hashedPassword string) error
	DeleteUser(user *table.UserTable) error
}
