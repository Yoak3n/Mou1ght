package interfaces

import "Mou1ght/internal/domain/model/table"

type AttachmentRepository interface {
	CreateAttachment(attachment *table.AttachmentTable) error
	GetAttachmentByID(id string) (*table.AttachmentTable, error)
	GetAttachmentsByIDs(ids []string) ([]table.AttachmentTable, error)
	// Unscoped：软删记录仍占用 StoragePath 唯一索引，重传同文件时必须能查到以便恢复
	GetAttachmentBySha256(sha256 string, size int64) (*table.AttachmentTable, error)
	RestoreAttachment(id string) error
	ListAttachments() ([]table.AttachmentTable, error)
	DeleteAttachment(id string) error
}
