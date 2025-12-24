package handler

import (
	"Mou1ght/consts"
	"Mou1ght/internal/api/controller"
	"Mou1ght/internal/pkg/util"
	util2 "Mou1ght/pkg/util"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetAttachmentList(c *fiber.Ctx) error {
	attachmentList := controller.GetAttachmentList()
	return util.SuccessResponse(c, attachmentList)
}

func UploadAttachment(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}

	// 兼容多文件上传
	// Get all files from "file" key:
	files := form.File["file"]
	// => []*multipart.FileHeader
	filesPaths := make([]string, len(files))
	filePath := ""
	// Loop through files:
	for index, file := range files {
		typ := strings.Split(file.Header.Get("Content-Type"), "/")
		if len(typ) > 1 {
			err := util2.CreateDirNotExists(consts.Upload + "/" + typ[1])
			if err != nil {
				continue
			}
			filePath = fmt.Sprintf("./%s/%s/%s", consts.Upload, typ[1], file.Filename)
		} else {
			filePath = fmt.Sprintf("./%s/%s", consts.Upload, file.Filename)
		}
		// Save the files to disk:
		err := c.SaveFile(file, filePath)
		// Check for errors
		if err != nil {
			return util.ErrorResponse(c, 500, err.Error())
		}
		// Return the route path to the file
		filesPaths[index], _ = strings.CutPrefix(filePath, fmt.Sprintf("./%s", consts.Cache))
	}
	return util.SuccessResponse(c, filesPaths)
}
