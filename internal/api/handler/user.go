package handler

import (
	"Mou1ght/internal/api/controller"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/pkg/util"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	loginRequest := &request.UserLoginRequest{}
	err := c.BodyParser(loginRequest)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	if len(loginRequest.Password) < 6 {
		return util.ErrorResponse(c, 400, errors.New("password length must be greater than 6").Error())
	}
	id, err := controller.UserLoginCheck(loginRequest)
	if err != nil {
		return util.ErrorResponse(c, 401, err.Error())
	}
	token, err := util.ReleaseToken(id)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, fiber.Map{"token": token})
}

func Register(c *fiber.Ctx) error {
	registerRequest := &request.UserRegisterRequest{}
	err := c.BodyParser(registerRequest)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	if len(registerRequest.Password) < 6 {
		return util.ErrorResponse(c, 400, errors.New("password length must be greater than 6").Error())
	}
	record, err := controller.UserRegisterCheck(registerRequest)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	token, err := util.ReleaseToken(record.ID)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, fiber.Map{"token": token, "name": record.UserName})
}

func Info(c *fiber.Ctx) error {
	userId := c.Locals("uid").(string)
	if userId == "" {
		return util.ErrorResponse(c, 401, "Unauthorized")
	}
	info, err := controller.UserInfo(userId)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, fiber.Map{"user": info})
}

func Logout(c *fiber.Ctx) error {
	err := util.ClearToken(c.Locals("token").(string))
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil)
}
