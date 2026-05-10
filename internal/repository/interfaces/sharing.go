package interfaces

import (
	"Mou1ght/internal/domain/model/table"
	"time"
)

type SharingRepository interface {
	CreateSharing(sharing *table.SharingTable) error
	UpdateSharing(sharing *table.SharingTable) error
	AddViewCountSharing(id string) error
	AddLikeCountSharing(id string) error
	GetSharingsByAuthorID(authorID string, desc bool) ([]table.SharingTable, error)
	GetSharings(startDate, endDate *time.Time) ([]*table.SharingTable, error)
	GetSharingByID(id string) (*table.SharingTable, error)
	DeleteSharingByID(id string) error
}
