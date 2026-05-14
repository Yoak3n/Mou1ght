package handler

import (
	"Mou1ght/internal/api/service"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/pkg/util"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

func (u *UserHandler) Login(c *fiber.Ctx) error {
	loginRequest := &request.UserLoginRequest{}
	err := c.BodyParser(loginRequest)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	if len(loginRequest.Password) < 6 {
		return util.ErrorResponse(c, 400, errors.New("password length must be greater than 6").Error())
	}
	id, err := u.userSvc.UserLoginCheck(loginRequest)
	if err != nil {
		return util.ErrorResponse(c, 401, err.Error())
	}
	token, err := util.ReleaseToken(id)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, fiber.Map{"token": token})
}

func (u *UserHandler) Register(c *fiber.Ctx) error {
	registerRequest := &request.UserRegisterRequest{}
	err := c.BodyParser(registerRequest)
	if err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	if len(registerRequest.Password) < 6 {
		return util.ErrorResponse(c, 400, errors.New("password length must be greater than 6").Error())
	}
	record, err := u.userSvc.UserRegisterCheck(registerRequest)
	if err != nil {
		if errors.Is(err, service.ErrRegistrationDisabled) {
			return util.ErrorResponse(c, 403, err.Error())
		}
		return util.ErrorResponse(c, 400, err.Error())
	}
	token, err := util.ReleaseToken(record.ID)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, fiber.Map{"token": token, "name": record.UserName})
}

func (u *UserHandler) Info(c *fiber.Ctx) error {
	userId := c.Locals("uid").(string)
	if userId == "" {
		return util.ErrorResponse(c, 401, "Unauthorized")
	}
	info, err := u.userSvc.UserInfo(userId)
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, fiber.Map{"user": info})
}

func (u *UserHandler) Logout(c *fiber.Ctx) error {
	err := util.ClearToken(c.Locals("token").(string))
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, nil)
}
