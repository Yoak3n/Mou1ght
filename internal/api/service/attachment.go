package service

import (
	"Mou1ght/internal/domain/entity"
	"strings"
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
