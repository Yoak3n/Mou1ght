package instance

import (
    "Mou1ght/internal/domain/model/table"
    "time"
    "gorm.io/gorm"
)

func (d *Database) GetSharingsByAuthorID(authorID string, desc bool) ([]table.SharingTable, error) {
	sharings := make([]table.SharingTable, 0)
	order := "created_at ASC"
	if desc {
		order = "created_at DESC"
	}
	err := d.DB.Where("author_id = ?", authorID).Order(order).Find(&sharings).Error
	if err != nil {
		return nil, err
	}
	return sharings, nil
}

func (d *Database) GetSharings(startDate, endDate *time.Time) ([]*table.SharingTable, error) {
	sharings := make([]*table.SharingTable, 0)
	if startDate != nil {

	}
	err := d.DB.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&sharings).Error
	if err != nil {
		return nil, err
	}
	return sharings, nil
}

func (d *Database) CreateSharing(sharing *table.SharingTable) error {
    return d.DB.Create(&sharing).Error
}

func (d *Database) UpdateSharing(sharing *table.SharingTable) error {
    return d.DB.Save(&sharing).Error
}

func (d *Database) AddViewCountSharing(id string) error {
    return d.DB.Where("id = ? AND status = 1", id).Update("view", gorm.Expr("view + 1")).Error
}

func (d *Database) AddLikeCountSharing(id string) error {
    return d.DB.Where("id = ? AND status = 1", id).Update("like", gorm.Expr("like + 1")).Error
}

func (d *Database) GetSharingByID(id string) (*table.SharingTable, error) {
    sharing := &table.SharingTable{}
    err := d.DB.Where("id = ?", id).First(&sharing).Error
    return sharing, err
}

func (d *Database) DeleteSharingByID(id string) error {
    return d.DB.Where("id = ?", id).Delete(&table.SharingTable{}).Error
}
