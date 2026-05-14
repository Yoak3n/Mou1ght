package entity

import (
	"path"
	"time"
)

type Attachment struct {
	FileName    string          `json:"file_name,omitempty"`
	FilePath    string          `json:"file_path,omitempty"`
	DirChildren []Attachment    `json:"dir_children,omitempty"`
	Info        *AttachmentInfo `json:"info,omitempty"`
}

type AttachmentInfo struct {
	Size    int64     `json:"size,omitempty"`
	ModTime time.Time `json:"mod_time"`
	IsDir   bool      `json:"is_dir,omitempty"`
}

type AttachmentEntity struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	OriginalName string `json:"original_name,omitempty"`
	Size         int64  `json:"size,omitempty"`
	Mime         string `json:"mime,omitempty"`
}

func NewAttachmentEntityFromPath(p string) AttachmentEntity {
	return AttachmentEntity{
		URL:          p,
		OriginalName: path.Base(p),
	}
}

func NewAttachmentsEntityFromPaths(paths []string) []AttachmentEntity {
	attachments := make([]AttachmentEntity, 0, len(paths))
	for _, p := range paths {
		attachments = append(attachments, NewAttachmentEntityFromPath(p))
	}
	return attachments
}
