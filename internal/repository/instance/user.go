package instance

import (
	"Mou1ght/internal/domain/model/table"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (d *UserRepository) CreateUser(user *table.UserTable) error {
	return d.db.Create(user).Error
}

func (d *UserRepository) GetUser(id string) (*table.UserTable, error) {
	user := &table.UserTable{}
	err := d.db.Where("id = ?", id).First(user).Error
	return user, err
}

func (d *UserRepository) GetUserByName(name string) (*table.UserTable, error) {
	user := &table.UserTable{}
	err := d.db.Where("user_name = ?", name).First(user).Error
	return user, err
}

func (d *UserRepository) QueryUsers(username []string) ([]table.UserTable, error) {
	users := make([]table.UserTable, 0)
	err := d.db.Where("user_name IN ?", username).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (d *UserRepository) CountUsers() (int64, error) {
	var count int64
	err := d.db.Model(&table.UserTable{}).Count(&count).Error
	return count, err
}

func (d *UserRepository) UpdateUser(user *table.UserTable) error {
	return d.db.Save(user).Error
}

func (d *UserRepository) UpdateUserProfile(id string, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}
	return d.db.Model(&table.UserTable{}).Where("id = ?", id).Updates(fields).Error
}

func (d *UserRepository) UpdateUserPassword(id, hashedPassword string) error {
	return d.db.Model(&table.UserTable{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

func (d *UserRepository) DeleteUser(user *table.UserTable) error {
	return d.db.Delete(user).Error
}
