package service

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/interfaces"
	"errors"
)

type SharingService struct {
	sharings     interfaces.SharingRepository
	tags         interfaces.TagRepository
	attachments  interfaces.AttachmentRepository
	sharingLinks interfaces.SharingAttachmentLinkRepository
}

func NewSharingService(sharings interfaces.SharingRepository, tags interfaces.TagRepository, attachments interfaces.AttachmentRepository, sharingLinks interfaces.SharingAttachmentLinkRepository) *SharingService {
	return &SharingService{sharings: sharings, tags: tags, attachments: attachments, sharingLinks: sharingLinks}
}

func (s *SharingService) CreateSharing(req *request.CreateSharingRequest) error {
	sid := util.GenSharingID()
	record := &table.SharingTable{
		PostBase: table.PostBase{
			ID:      sid,
			Content: req.Content,
		},
		AuthorID:   req.Author,
		Attachment: req.Attachment,
	}
	if err := s.validateAttachmentIDs(req.AttachmentIDs); err != nil {
		return err
	}
	err := s.sharings.CreateSharing(record)
	if err != nil {
		return err
	}
	err = s.sharingLinks.ReplaceSharingAttachments(sid, req.AttachmentIDs)
	if err != nil {
		return err
	}
	tagIDs := make([]string, len(req.Tags))
	for i, tag := range req.Tags {
		tagIDs[i] = tag.ID
	}
	err = s.CreateTagsLinkToSharing(tagIDs, sid)
	if err != nil {
		return err
	}
	return nil
}

func (s *SharingService) UpdateSharing(req *request.UpdateSharingRequest) error {
	record := &table.SharingTable{
		PostBase: table.PostBase{
			ID:      req.ID,
			Content: req.Content,
		},
		AuthorID:   req.Author,
		Attachment: req.Attachment,
	}
	if err := s.validateAttachmentIDs(req.AttachmentIDs); err != nil {
		return err
	}
	err := s.sharings.UpdateSharing(record)
	if err != nil {
		return err
	}
	err = s.sharingLinks.ReplaceSharingAttachments(req.ID, req.AttachmentIDs)
	if err != nil {
		return err
	}
	tagsIDs := make(map[string]bool)
	for _, tag := range req.Tags {
		tagsIDs[tag.ID] = true
	}
	err = s.tags.UpdateTargetLinks(req.ID, 2, tagsIDs)
	if err != nil {
		return err
	}
	return nil
}

func (s *SharingService) ViewSharing(id string) error {
	return s.sharings.AddViewCountSharing(id)
}

func (s *SharingService) LikeSharing(id string) error {
	return s.sharings.AddLikeCountSharing(id)
}

func (s *SharingService) GetSharingByID(id string) (*table.SharingTable, error) {
	record, err := s.sharings.GetSharingByID(id)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (s *SharingService) DeleteSharingByID(id string) error {
	err := s.sharings.DeleteSharingByID(id)
	if err != nil {
		return err
	}
	_ = s.sharingLinks.DeleteBySharingID(id)
	err = s.tags.DeleteTagLinkFromTarget(id, 2)
	if err != nil {
		return err
	}
	return nil
}

func (s *SharingService) validateAttachmentIDs(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	records, err := s.attachments.GetAttachmentsByIDs(ids)
	if err != nil {
		return err
	}
	rm := make(map[string]bool, len(records))
	for _, r := range records {
		rm[r.ID] = true
	}
	for _, id := range ids {
		if id == "" {
			continue
		}
		if _, ok := rm[id]; !ok {
			return errors.New("attachment not exist")
		}
	}
	return nil
}

func (s *SharingService) CreateTagsLinkToSharing(tags []string, sharingID string) error {
	for _, tag := range tags {
		lid := util.GenTagLinkID()
		record := &table.TagLinkTable{
			ID:         lid,
			TargetID:   sharingID,
			TargetType: 2,
			TagID:      tag,
		}
		err := s.tags.CreateTagLink(record)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SharingService) GetSharingsByAuthorID(authorID string, descend bool) ([]table.SharingTable, error) {
	return s.sharings.GetSharingsByAuthorID(authorID, descend)
}
