package instance

import (
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"

	"gorm.io/gorm"
)

type AttachmentLinkRepository struct {
	db *gorm.DB
}

func NewAttachmentLinkRepository(db *gorm.DB) *AttachmentLinkRepository {
	return &AttachmentLinkRepository{db: db}
}

func (r *AttachmentLinkRepository) ReplaceSharingAttachments(sharingID string, attachmentIDs []string) error {
	err := r.DeleteBySharingID(sharingID)
	if err != nil {
		return err
	}
	if len(attachmentIDs) == 0 {
		return nil
	}
	for i, id := range attachmentIDs {
		link := &table.AttachmentLinkTable{
			ID:           util.GenAttachmentLinkID(),
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

func (r *AttachmentLinkRepository) DeleteBySharingID(sharingID string) error {
	return r.db.Where("sharing_id = ?", sharingID).Delete(&table.AttachmentLinkTable{}).Error
}

func (r *AttachmentLinkRepository) GetAttachmentIDsBySharingID(sharingID string) ([]string, error) {
	ids := make([]string, 0)
	err := r.db.Where("sharing_id = ?", sharingID).Order("sort ASC").Model(&table.AttachmentLinkTable{}).Pluck("attachment_id", &ids).Error
	return ids, err
}
