package instance

import (
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/repository/interfaces"
	"log"
	"time"

	"gorm.io/gorm"
)

type SharingRepository struct {
	db      *gorm.DB
	counter interfaces.PostCounter
}

func NewSharingRepository(db *gorm.DB, counter interfaces.PostCounter) *SharingRepository {
	return &SharingRepository{
		db:      db,
		counter: counter,
	}
}

func (s *SharingRepository) GetSharingsByAuthorID(authorID string, desc bool) ([]table.SharingTable, error) {
	sharings := make([]table.SharingTable, 0)
	order := "created_at ASC"
	if desc {
		order = "created_at DESC"
	}
	err := s.db.Model(&table.SharingTable{}).Where("author_id = ?", authorID).Order(order).Find(&sharings).Error
	if err != nil {
		return nil, err
	}
	return sharings, nil
}

func (s *SharingRepository) GetSharings(startDate, endDate *time.Time) ([]*table.SharingTable, error) {
	sharings := make([]*table.SharingTable, 0)
	var query *gorm.DB
	if startDate != nil {
		if endDate == nil {
			query = s.db.Model(&table.SharingTable{}).Where("created_at >= ?", startDate)
		} else {
			query = s.db.Model(&table.SharingTable{}).Where("created_at BETWEEN ? AND ?", startDate, endDate)
		}
	} else {
		if endDate == nil {
			query = s.db.Model(&table.SharingTable{})
		} else {
			query = s.db.Model(&table.SharingTable{}).Where("created_at <= ?", endDate)
		}
	}
	err := query.Order("created_at DESC").Find(&sharings).Error
	if err != nil {
		return nil, err
	}
	log.Println(len(sharings))
	return sharings, nil
}

func (s *SharingRepository) CreateSharing(sharing *table.SharingTable) error {
	return s.db.Create(&sharing).Error
}

func (s *SharingRepository) UpdateSharing(sharing *table.SharingTable) error {
	return s.db.Save(&sharing).Error
}

func (s *SharingRepository) AddViewCountSharing(id string) error {
	s.counter.BumpView("sharing", id, 1)
	return nil
}

func (s *SharingRepository) AddLikeCountSharing(id string) error {
	s.counter.BumpLike("sharing", id, 1)
	return nil
}

func (s *SharingRepository) GetSharingByID(id string) (*table.SharingTable, error) {
	sharing := &table.SharingTable{}
	err := s.db.Model(&table.SharingTable{}).Where("id = ?", id).First(&sharing).Error
	return sharing, err
}

func (s *SharingRepository) DeleteSharingByID(id string) error {
	return s.db.Where("id = ?", id).Delete(&table.SharingTable{}).Error
}
