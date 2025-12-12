package handler

import (
	"Mou1ght/internal/api/controller"
	"Mou1ght/internal/domain/model/schema/console"
	"Mou1ght/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
)

func GetAllSetting(c *fiber.Ctx) error {
	setting, err := controller.GetAllSetting()
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, setting)
}

func GetBlogSetting(c *fiber.Ctx) error {
	setting, err := controller.GetBlogSetting()
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, setting)
}

func UpdateBlogSetting(c *fiber.Ctx) error {
	setting := new(console.BlogSetting)
	if err := c.BodyParser(setting); err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	if err := controller.UpdateBlogSetting(*setting); err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil)
}
