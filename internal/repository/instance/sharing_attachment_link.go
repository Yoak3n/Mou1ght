package instance

import (
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"

	"gorm.io/gorm"
)

type SharingAttachmentLinkRepository struct {
	db *gorm.DB
}

func NewSharingAttachmentLinkRepository(db *gorm.DB) *SharingAttachmentLinkRepository {
	return &SharingAttachmentLinkRepository{db: db}
}

func (r *SharingAttachmentLinkRepository) ReplaceSharingAttachments(sharingID string, attachmentIDs []string) error {
	err := r.DeleteBySharingID(sharingID)
	if err != nil {
		return err
	}
	if len(attachmentIDs) == 0 {
		return nil
	}
	for i, id := range attachmentIDs {
		link := &table.SharingAttachmentLinkTable{
			ID:           util.GenSharingAttachmentLinkID(),
			SharingID:    sharingID,
			AttachmentID: id,
			Sort:         i,
		}
		if e := r.db.Create(link).Error; e != nil {
			err = e
		}
	}
	return err
}

func (r *SharingAttachmentLinkRepository) DeleteBySharingID(sharingID string) error {
	return r.db.Where("sharing_id = ?", sharingID).Delete(&table.SharingAttachmentLinkTable{}).Error
}

func (r *SharingAttachmentLinkRepository) GetAttachmentIDsBySharingID(sharingID string) ([]string, error) {
	ids := make([]string, 0)
	err := r.db.Where("sharing_id = ?", sharingID).Order("sort ASC").Model(&table.SharingAttachmentLinkTable{}).Pluck("attachment_id", &ids).Error
	return ids, err
}
