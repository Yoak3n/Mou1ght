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
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	Size     int64  `json:"size"`
	IsDir    bool   `json:"is_dir"`
}

func NewAttachmentEntityFromPath(p string) AttachmentEntity {
	return AttachmentEntity{
		FileName: path.Base(p),
		FilePath: p,
	}
}

func NewAttachmentsEntityFromPaths(paths []string) []AttachmentEntity {
	attachments := make([]AttachmentEntity, 0, len(paths))
	for _, p := range paths {
		attachments = append(attachments, NewAttachmentEntityFromPath(p))
	}
	return attachments
}
