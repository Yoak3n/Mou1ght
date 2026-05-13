package service

import (
	"Mou1ght/consts"
	"Mou1ght/internal/domain/entity"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/interfaces"
	fileutil "Mou1ght/pkg/util"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

type AttachmentService struct {
	attachments interfaces.AttachmentRepository
}

func NewAttachmentService(attachments interfaces.AttachmentRepository) *AttachmentService {
	return &AttachmentService{attachments: attachments}
}

func (s *AttachmentService) Upload(files []*multipart.FileHeader) ([]entity.AttachmentEntity, error) {
	if len(files) == 0 {
		return []entity.AttachmentEntity{}, nil
	}
	entities := make([]entity.AttachmentEntity, 0, len(files))
	for _, f := range files {
		if f == nil {
			continue
		}
		e, err := s.uploadOne(f)
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}
	return entities, nil
}

func (s *AttachmentService) uploadOne(file *multipart.FileHeader) (entity.AttachmentEntity, error) {
	if file.Size <= 0 {
		return entity.AttachmentEntity{}, errors.New("empty file")
	}
	if file.Size > 50*1024*1024 {
		return entity.AttachmentEntity{}, errors.New("file too large")
	}

	mimeType := file.Header.Get("Content-Type")
	mimeType = strings.TrimSpace(mimeType)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	hash := sha256.New()
	src, err := file.Open()
	if err != nil {
		return entity.AttachmentEntity{}, err
	}
	_, err = io.Copy(hash, src)
	_ = src.Close()
	if err != nil {
		return entity.AttachmentEntity{}, err
	}
	shaHex := hex.EncodeToString(hash.Sum(nil))

	existing, err := s.attachments.GetAttachmentBySha256(shaHex, file.Size)
	if err == nil && existing != nil && existing.ID != "" {
		return attachmentEntityFromTable(existing), nil
	}

	ext := resolveExt(file.Filename, mimeType)
	dir := resolveDir(mimeType)
	now := time.Now()
	storagePath := path.Join(dir, now.Format("2006"), now.Format("01"), shaHex+ext)

	dstDir := path.Join(consts.Upload, dir, now.Format("2006"), now.Format("01"))
	if e := fileutil.CreateDirNotExists(dstDir); e != nil {
		return entity.AttachmentEntity{}, e
	}

	dstFullPath := path.Join(consts.Upload, storagePath)
	src2, err := file.Open()
	if err != nil {
		return entity.AttachmentEntity{}, err
	}
	dst, err := os.OpenFile(dstFullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		_ = src2.Close()
		return entity.AttachmentEntity{}, err
	}
	_, err = io.Copy(dst, src2)
	_ = dst.Close()
	_ = src2.Close()
	if err != nil {
		_ = os.Remove(dstFullPath)
		return entity.AttachmentEntity{}, err
	}

	record := &table.AttachmentTable{
		ID:           util.GenAttachmentID(),
		OriginalName: path.Base(file.Filename),
		StoragePath:  storagePath,
		Mime:         mimeType,
		Sha256:       shaHex,
		Size:         file.Size,
	}
	if err := s.attachments.CreateAttachment(record); err != nil {
		_ = os.Remove(dstFullPath)
		return entity.AttachmentEntity{}, err
	}
	return attachmentEntityFromTable(record), nil
}

func attachmentEntityFromTable(t *table.AttachmentTable) entity.AttachmentEntity {
	if t == nil {
		return entity.AttachmentEntity{}
	}
	return entity.AttachmentEntity{
		ID:           t.ID,
		URL:          "/upload/" + strings.TrimPrefix(t.StoragePath, "/"),
		OriginalName: t.OriginalName,
		Size:         t.Size,
		Mime:         t.Mime,
	}
}

func resolveDir(mimeType string) string {
	if i := strings.Index(mimeType, "/"); i > 0 {
		return mimeType[:i]
	}
	return "file"
}

func resolveExt(filename string, mimeType string) string {
	ext := strings.ToLower(path.Ext(filename))
	if ext != "" && len(ext) <= 10 {
		return ext
	}
	if exts, err := mime.ExtensionsByType(mimeType); err == nil && len(exts) > 0 {
		return exts[0]
	}
	return ""
}
