package interfaces

type SharingAttachmentLinkRepository interface {
	ReplaceSharingAttachments(sharingID string, attachmentIDs []string) error
	DeleteBySharingID(sharingID string) error
	GetAttachmentIDsBySharingID(sharingID string) ([]string, error)
}
