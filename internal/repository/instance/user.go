package instance

import "Mou1ght/internal/domain/model/table"

func (d *Database) CreateUser(user *table.UserTable) error {
	return d.DB.Create(user).Error
}

func (d *Database) GetUser(id string) (*table.UserTable, error) {
	user := &table.UserTable{}
	err := d.DB.First(user).Where("id = ?", id).Error
	return user, err
}

func (d *Database) GetUserByName(name string) (*table.UserTable, error) {
	user := &table.UserTable{}
	err := d.DB.First(user).Where("user_name = ?", name).Error
	return user, err
}

func (d *Database) QueryUsers(username []string) ([]table.UserTable, error) {
	users := make([]table.UserTable, 0)
	err := d.DB.Find(&users).Where("user_name IN ?", username).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (d *Database) UpdateUser(user *table.UserTable) error {
	return d.DB.Save(user).Error
}
func (d *Database) DeleteUser(user *table.UserTable) error {
	return d.DB.Delete(user).Error
}
