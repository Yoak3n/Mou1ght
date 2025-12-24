package controller

import (
	"Mou1ght/internal/domain/entity"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/instance"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserLoginCheck(req *request.UserLoginRequest) (string, error) {
	user, err := instance.UseDatabase().GetUserByName(req.UserName)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("password incorrect")
	}
	user.LastLogin = time.Now()
	_ = instance.UseDatabase().UpdateUser(user)
	return user.ID, nil
}

func UserRegisterCheck(req *request.UserRegisterRequest) (*table.UserTable, error) {
	uid := ""
	for {
		uid = util.GenUserID()
		_, err := instance.UseDatabase().GetUser(uid)
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
	err = instance.UseDatabase().CreateUser(record)
	if err != nil {
		// 一般是因为用户名重复
		return nil, errors.New("user name perhaps already exists")
	}
	return record, nil
}

func UserInfo(uid string) (*entity.UserEntity, error) {
	record, err := instance.UseDatabase().GetUser(uid)
	if err != nil {
		return nil, err
	}
	return entity.NewUserEntityFromTable(record, true), nil
}

func AuthorListWithPost(req *request.PostListRequest) map[string]any {
	ret := make(map[string]any)
	filter := req.Filter
	descend := filter.Sort == "desc"
	authors, err := instance.UseDatabase().QueryUsers(req.Data.Keyword)
	if err != nil {
		return nil
	}
	es := make([]*entity.UserWithPostEntity, 0)
	for _, author := range authors {
		articles, _ := instance.UseDatabase().GetArticlesByAuthorID(author.ID, descend)
		sharings, _ := instance.UseDatabase().GetSharingsByAuthorID(author.ID, descend)
		e := entity.NewUserWithPostEntityFromTable(&author, sharings, articles)
		es = append(es, e)
	}
	ret["authors"] = es
	return nil
}
