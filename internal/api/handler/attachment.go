package handler

import (
	"Mou1ght/internal/api/service"
	"Mou1ght/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
)

type AttachmentHandler struct {
	attachmentService *service.AttachmentService
}

func NewAttachmentHandler(attachmentService *service.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{attachmentService: attachmentService}
}

func (h *AttachmentHandler) GetAttachmentList(c *fiber.Ctx) error {
	attachments, err := h.attachmentService.ListAll()
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, fiber.Map{
		"attachments": attachments,
	}, "")
}

func (h *AttachmentHandler) UploadAttachment(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}

	files := form.File["file"]
	attachments, err := h.attachmentService.Upload(files)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, fiber.Map{
		"attachments": attachments,
	}, "")
}
