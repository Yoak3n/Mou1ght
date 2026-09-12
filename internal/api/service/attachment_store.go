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
	links       interfaces.SharingAttachmentLinkRepository
}

func NewAttachmentService(attachments interfaces.AttachmentRepository, links interfaces.SharingAttachmentLinkRepository) *AttachmentService {
	return &AttachmentService{attachments: attachments, links: links}
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

	ext := resolveExt(file.Filename, mimeType)
	dir := resolveDir(mimeType)
	now := time.Now()
	storagePath := path.Join(dir, now.Format("2006"), now.Format("01"), shaHex+ext)

	existing, err := s.attachments.GetAttachmentBySha256(shaHex, file.Size)
	if err == nil && existing != nil && existing.ID != "" {
		// 同内容已存在：可能是软删残留（uniqueIndex 仍占位），恢复并确保磁盘文件在
		if e := writeAttachmentFile(file, existing.StoragePath); e != nil {
			return entity.AttachmentEntity{}, e
		}
		if existing.DeletedAt.Valid {
			if e := s.attachments.RestoreAttachment(existing.ID); e != nil {
				return entity.AttachmentEntity{}, e
			}
		}
		return attachmentEntityFromTable(existing), nil
	}

	// 新文件：路径用已有记录的 storagePath 规则；避免与软删行撞 unique 时先写盘再入库
	if e := writeAttachmentFile(file, storagePath); e != nil {
		return entity.AttachmentEntity{}, e
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
		_ = os.Remove(path.Join(consts.Upload, strings.TrimPrefix(storagePath, "/")))
		return entity.AttachmentEntity{}, err
	}
	return attachmentEntityFromTable(record), nil
}

func writeAttachmentFile(file *multipart.FileHeader, storagePath string) error {
	dstDir := path.Join(consts.Upload, path.Dir(strings.TrimPrefix(storagePath, "/")))
	if e := fileutil.CreateDirNotExists(dstDir); e != nil {
		return e
	}
	dstFullPath := path.Join(consts.Upload, strings.TrimPrefix(storagePath, "/"))
	src, err := file.Open()
	if err != nil {
		return err
	}
	dst, err := os.OpenFile(dstFullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		_ = src.Close()
		return err
	}
	_, err = io.Copy(dst, src)
	_ = dst.Close()
	_ = src.Close()
	if err != nil {
		_ = os.Remove(dstFullPath)
		return err
	}
	return nil
}

func attachmentEntityFromTable(t *table.AttachmentTable) entity.AttachmentEntity {
	if t == nil {
		return entity.AttachmentEntity{}
	}
	rel := "/upload/" + strings.TrimPrefix(t.StoragePath, "/")
	return entity.AttachmentEntity{
		ID:           t.ID,
		URL:          rel,
		FilePath:     rel,
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
