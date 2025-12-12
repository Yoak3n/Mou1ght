package handler

import (
    "Mou1ght/internal/api/controller"
    "Mou1ght/internal/domain/model/schema/request"
    "Mou1ght/internal/pkg/util"

    "github.com/gofiber/fiber/v2"
)

func CreateSharing(c *fiber.Ctx) error {
    req := &request.CreateSharingRequest{}
    err := c.BodyParser(req)
    if err != nil {
        return util.ErrorResponse(c, 400, err.Error())
    }
    err = controller.CreateSharing(req)
    if err != nil {
        return util.ErrorResponse(c, 500, err.Error())
    }
    return util.SuccessResponse(c, nil, "Create sharing successfully")
}

func DeleteSharing(c *fiber.Ctx) error {
    id := c.Params("id")
    if id == "" {
        return util.ErrorResponse(c, 400, "id is required")
    }
    err := controller.DeleteSharingByID(id)
    if err != nil {
        return util.ErrorResponse(c, 500, err.Error())
    }
    return util.SuccessResponse(c, nil, "Delete sharing successfully")
}

func UpdateSharing(c *fiber.Ctx) error {
    req := &request.UpdateSharingRequest{}
    err := c.BodyParser(req)
    if err != nil {
        return util.ErrorResponse(c, 400, err.Error())
    }
    err = controller.UpdateSharing(req)
    if err != nil {
        return util.ErrorResponse(c, 500, err.Error())
    }
    return util.SuccessResponse(c, nil, "Update sharing successfully")
}

func DetailSharing(c *fiber.Ctx) error {
    id := c.Params("id")
    if id == "" {
        return util.ErrorResponse(c, 400, "id is required")
    }
    sharing, err := controller.GetSharingByID(id)
    if err != nil {
        return util.ErrorResponse(c, 500, err.Error())
    }
    return util.SuccessResponse(c, sharing)
}

