package interfaces

// SharingAttachmentLinkRepository 管理说说与文章的附件链接（AttachmentLinkTable 同时承载两类目标）
type SharingAttachmentLinkRepository interface {
	ReplaceSharingAttachments(sharingID string, attachmentIDs []string) error
	DeleteBySharingID(sharingID string) error
	GetAttachmentIDsBySharingID(sharingID string) ([]string, error)
	ReplaceArticleAttachments(articleID string, attachmentIDs []string) error
	DeleteByArticleID(articleID string) error
	GetAttachmentIDsByArticleID(articleID string) ([]string, error)
	CountByAttachmentID(attachmentID string) (int64, error)
}
