package instance

import (
	"Mou1ght/internal/domain/model/table"

	"gorm.io/gorm"
)

type AttachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{db: db}
}

func (r *AttachmentRepository) CreateAttachment(attachment *table.AttachmentTable) error {
	return r.db.Create(attachment).Error
}

func (r *AttachmentRepository) GetAttachmentByID(id string) (*table.AttachmentTable, error) {
	attachment := &table.AttachmentTable{}
	err := r.db.Where("id = ?", id).First(attachment).Error
	return attachment, err
}

func (r *AttachmentRepository) GetAttachmentsByIDs(ids []string) ([]table.AttachmentTable, error) {
	if len(ids) == 0 {
		return []table.AttachmentTable{}, nil
	}
	attachments := make([]table.AttachmentTable, 0, len(ids))
	err := r.db.Where("id IN ?", ids).Find(&attachments).Error
	return attachments, err
}

func (r *AttachmentRepository) GetAttachmentBySha256(sha256 string, size int64) (*table.AttachmentTable, error) {
	attachment := &table.AttachmentTable{}
	result := r.db.Where("sha256 = ? AND size = ?", sha256, size).First(attachment)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return attachment, result.Error
}

func (r *AttachmentRepository) ListAttachments() ([]table.AttachmentTable, error) {
	attachments := make([]table.AttachmentTable, 0)
	err := r.db.Order("created_at DESC").Find(&attachments).Error
	return attachments, err
}
