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

func (s *UserService) AuthorListWithPost(req *request.PostListRequest) []table.UserTable {
	authors, err := s.users.QueryUsers(req.Data.Keyword)
	if err != nil {
		return nil
	}

	return authors
}
