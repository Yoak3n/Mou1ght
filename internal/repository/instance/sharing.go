package instance

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/repository/interfaces"

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

func (s *SharingRepository) GetSharings(opts request.ListOptions) ([]*table.SharingTable, int64, error) {
	sharings := make([]*table.SharingTable, 0)
	query := s.db.Model(&table.SharingTable{})
	switch {
	case opts.StartDate != nil && opts.EndDate != nil:
		query = query.Where("created_at BETWEEN ? AND ?", opts.StartDate, opts.EndDate)
	case opts.StartDate != nil:
		query = query.Where("created_at >= ?", opts.StartDate)
	case opts.EndDate != nil:
		query = query.Where("created_at <= ?", opts.EndDate)
	}
	if opts.OnlyPublished {
		query = query.Where("status = ?", 1)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "created_at DESC"
	if opts.AscendingOrder {
		order = "created_at ASC"
	}
	query = query.Order(order)
	if opts.Paginated() {
		query = query.Offset(opts.Offset()).Limit(opts.PageSize)
	}
	if err := query.Find(&sharings).Error; err != nil {
		return nil, 0, err
	}
	return sharings, total, nil
}

func (s *SharingRepository) CreateSharing(sharing *table.SharingTable) error {
	return s.db.Create(&sharing).Error
}

func (s *SharingRepository) UpdateSharing(sharing *table.SharingTable) error {
	return s.db.Model(&table.SharingTable{}).
		Where("id = ?", sharing.ID).
		Updates(map[string]any{
			"content":   sharing.Content,
			"author_id": sharing.AuthorID,
			"status":    sharing.Status,
		}).Error
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
