package controller

import (
	"Mou1ght/consts"
	"Mou1ght/internal/domain/entity"
	"os"
	"path"
	"strings"
)

func GetAttachmentList() entity.Attachment {
	attachmentList := lookThroughAttachment(consts.Upload)
	return attachmentList
}

func lookThroughAttachment(entry string) entity.Attachment {
	p, _ := strings.CutPrefix(entry, consts.Cache+"/")
	ea := entity.Attachment{
		FileName:    path.Base(p),
		FilePath:    p,
		DirChildren: make([]entity.Attachment, 0),
		Info: &entity.AttachmentInfo{
			IsDir: true,
		},
	}
	dir, err := os.ReadDir(entry)
	if err != nil {
		return ea
	}
	for _, sub := range dir {
		if sub.IsDir() {
			subDir := lookThroughAttachment(entry + "/" + sub.Name())
			ea.DirChildren = append(ea.DirChildren, subDir)
		}
		info, e := sub.Info()
		if e != nil {
			return entity.Attachment{}
		}
		eac := entity.Attachment{
			FileName: sub.Name(),
			FilePath: p + "/" + sub.Name(),
			Info: &entity.AttachmentInfo{
				Size:    info.Size(),
				ModTime: info.ModTime(),
				IsDir:   sub.IsDir(),
			},
		}
		ea.DirChildren = append(ea.DirChildren, eac)
	}
	return ea
}
