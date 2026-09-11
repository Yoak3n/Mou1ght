package service

import (
	"Mou1ght/consts"
	"Mou1ght/internal/domain/entity"
	"errors"
	"os"
	"path"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrAttachmentNotFound   = errors.New("attachment not found")
	ErrAttachmentInUse      = errors.New("attachment is referenced by sharing")
	ErrAttachmentIDRequired = errors.New("id is required")
)

func (s *AttachmentService) ListAll() ([]entity.AttachmentEntity, error) {
	records, err := s.attachments.ListAttachments()
	if err != nil {
		return nil, err
	}
	entities := make([]entity.AttachmentEntity, 0, len(records))
	for i := range records {
		entities = append(entities, entity.AttachmentEntity{
			ID:           records[i].ID,
			URL:          "/upload/" + strings.TrimPrefix(records[i].StoragePath, "/"),
			OriginalName: records[i].OriginalName,
			Size:         records[i].Size,
			Mime:         records[i].Mime,
		})
	}
	return entities, nil
}

func (s *AttachmentService) Delete(id string) error {
	if id == "" {
		return ErrAttachmentIDRequired
	}
	record, err := s.attachments.GetAttachmentByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAttachmentNotFound
		}
		return err
	}
	if record == nil || record.ID == "" {
		return ErrAttachmentNotFound
	}

	used, err := s.links.CountByAttachmentID(id)
	if err != nil {
		return err
	}
	if used > 0 {
		return ErrAttachmentInUse
	}

	if err := s.attachments.DeleteAttachment(id); err != nil {
		return err
	}

	fullPath := path.Join(consts.Upload, strings.TrimPrefix(record.StoragePath, "/"))
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
