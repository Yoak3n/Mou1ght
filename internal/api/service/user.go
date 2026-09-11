package service

import (
	"Mou1ght/internal/domain/entity"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/interfaces"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrRegistrationDisabled = errors.New("registration disabled")

type UserService struct {
	users    interfaces.UserRepository
	articles interfaces.ArticleRepository
	sharings interfaces.SharingRepository
}

func NewUserService(users interfaces.UserRepository, articles interfaces.ArticleRepository, sharings interfaces.SharingRepository) *UserService {
	return &UserService{users: users, articles: articles, sharings: sharings}
}

func (s *UserService) UserLoginCheck(req *request.UserLoginRequest) (string, error) {
	user, err := s.users.GetUserByName(req.UserName)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("password incorrect")
	}
	user.LastLogin = time.Now()
	_ = s.users.UpdateUser(user)
	return user.ID, nil
}

func (s *UserService) UserRegisterCheck(req *request.UserRegisterRequest) (*table.UserTable, error) {
	count, err := s.users.CountUsers()
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrRegistrationDisabled
	}

	uid := ""
	for {
		uid = util.GenUserID()
		_, err := s.users.GetUser(uid)
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			break
		}
	}

	record := &table.UserTable{
		ID:       uid,
		UserName: req.UserName,
		Email:    req.Email,
		Role:     0,
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	record.Password = string(hashedPassword)
	err = s.users.CreateUser(record)
	if err != nil {
		// 一般是因为用户名重复
		return nil, errors.New("user name perhaps already exists")
	}
	return record, nil
}

func (s *UserService) UserInfo(uid string) (*entity.UserEntity, error) {
	record, err := s.users.GetUser(uid)
	if err != nil {
		return nil, err
	}
	return entity.NewUserEntityFromTable(record, true), nil
}

func (s *UserService) UpdateProfile(uid string, req *request.UpdateUserProfileRequest) (*entity.UserEntity, error) {
	if uid == "" {
		return nil, errors.New("uid is required")
	}
	user, err := s.users.GetUser(uid)
	if err != nil {
		return nil, err
	}

	fields := make(map[string]any)
	if req.UserName != "" && req.UserName != user.UserName {
		if _, err := s.users.GetUserByName(req.UserName); err == nil {
			return nil, errors.New("user name already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fields["user_name"] = req.UserName
	}
	if req.Email != user.Email {
		fields["email"] = req.Email
	}
	if req.Phone != user.Phone {
		fields["phone"] = req.Phone
	}
	if req.Avatar != user.Avatar {
		fields["avatar"] = req.Avatar
	}
	if len(fields) == 0 {
		return entity.NewUserEntityFromTable(user, true), nil
	}
	if err := s.users.UpdateUserProfile(uid, fields); err != nil {
		return nil, err
	}
	updated, err := s.users.GetUser(uid)
	if err != nil {
		return nil, err
	}
	return entity.NewUserEntityFromTable(updated, true), nil
}

func (s *UserService) ChangePassword(uid string, req *request.ChangePasswordRequest) error {
	if uid == "" {
		return errors.New("uid is required")
	}
	user, err := s.users.GetUser(uid)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("old password incorrect")
	}
	if req.NewPassword == req.OldPassword {
		return errors.New("new password must differ from old password")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.users.UpdateUserPassword(uid, string(hashed))
}

func (s *UserService) AuthorListWithPost(req *request.PostListRequest) []table.UserTable {
	authors, err := s.users.QueryUsers(req.Data.Keyword)
	if err != nil {
		return nil
	}

	return authors
}
