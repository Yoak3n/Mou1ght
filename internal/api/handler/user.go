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

// RegisterStatus 返回注册开关状态，供前台决定是否展示注册入口。
func (u *UserHandler) RegisterStatus(c *fiber.Ctx) error {
	open, err := u.userSvc.IsRegistrationOpen()
	if err != nil {
		return util.ErrorResponse(c, 500, err.Error())
	}
	return util.SuccessResponse(c, fiber.Map{"open": open})
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

func (u *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userId := c.Locals("uid").(string)
	if userId == "" {
		return util.ErrorResponse(c, 401, "Unauthorized")
	}
	req := &request.UpdateUserProfileRequest{}
	if err := c.BodyParser(req); err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	info, err := u.userSvc.UpdateProfile(userId, req)
	if err != nil {
		msg := err.Error()
		if msg == "user name already exists" {
			return util.ErrorResponse(c, 409, msg)
		}
		return util.ErrorResponse(c, 400, msg)
	}
	return util.SuccessResponse(c, fiber.Map{"user": info})
}

func (u *UserHandler) ChangePassword(c *fiber.Ctx) error {
	userId := c.Locals("uid").(string)
	if userId == "" {
		return util.ErrorResponse(c, 401, "Unauthorized")
	}
	req := &request.ChangePasswordRequest{}
	if err := c.BodyParser(req); err != nil {
		return util.ErrorResponse(c, 400, err.Error())
	}
	if len(req.NewPassword) < 6 {
		return util.ErrorResponse(c, 400, "password length must be greater than 6")
	}
	if req.OldPassword == "" {
		return util.ErrorResponse(c, 400, "old_password is required")
	}
	if err := u.userSvc.ChangePassword(userId, req); err != nil {
		msg := err.Error()
		if msg == "old password incorrect" {
			return util.ErrorResponse(c, 401, msg)
		}
		return util.ErrorResponse(c, 400, msg)
	}
	return util.SuccessResponse(c, nil, "Change password successfully")
}
