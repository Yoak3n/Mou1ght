package handler

import (
	"Mou1ght/internal/api/service"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
)

type SharingHandler struct {
	sharingService *service.SharingService
	dtoService     *service.DTOService
}

func NewSharingHandler(sharingService *service.SharingService, dtoService *service.DTOService) *SharingHandler {
	return &SharingHandler{sharingService: sharingService, dtoService: dtoService}
}

func (s *SharingHandler) CreateSharing(c *fiber.Ctx) error {
	req := &request.CreateSharingRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	err = s.sharingService.CreateSharing(req)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Create sharing successfully")
}

func (s *SharingHandler) DeleteSharing(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	err := s.sharingService.DeleteSharingByID(id)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Delete sharing successfully")
}

func (s *SharingHandler) UpdateSharing(c *fiber.Ctx) error {
	req := &request.UpdateSharingRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	err = s.sharingService.UpdateSharing(req)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil, "Update sharing successfully")
}

func (s *SharingHandler) DetailSharing(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.ErrorResponse(c, 400, "id is required")
	}
	sharing, err := s.sharingService.GetSharingByID(id)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	d := s.dtoService.GetSharingEntityFromTable(sharing)
	if d == nil {
		return util.ErrorResponse(c, 400, "sharing not exist")
	}
	return util.SuccessResponse(c, d)
}
