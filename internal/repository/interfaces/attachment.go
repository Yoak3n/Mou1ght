package interfaces

import "Mou1ght/internal/domain/model/table"

type AttachmentRepository interface {
	CreateAttachment(attachment *table.AttachmentTable) error
	GetAttachmentByID(id string) (*table.AttachmentTable, error)
	GetAttachmentsByIDs(ids []string) ([]table.AttachmentTable, error)
	GetAttachmentBySha256(sha256 string, size int64) (*table.AttachmentTable, error)
	ListAttachments() ([]table.AttachmentTable, error)
}
