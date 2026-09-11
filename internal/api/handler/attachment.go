package handler

import (
	"Mou1ght/internal/api/service"
	"Mou1ght/internal/pkg/util"
	"errors"

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

func (h *AttachmentHandler) DeleteAttachment(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, service.ErrAttachmentIDRequired.Error())
	}
	if err := h.attachmentService.Delete(id); err != nil {
		switch {
		case errors.Is(err, service.ErrAttachmentNotFound):
			return util.ErrorResponse(c, 404, err.Error())
		case errors.Is(err, service.ErrAttachmentInUse):
			return util.ErrorResponse(c, 409, err.Error())
		default:
			return util.ErrorResponse(c, 500, err.Error())
		}
	}
	return util.SuccessResponse(c, nil, "Delete attachment successfully")
}
