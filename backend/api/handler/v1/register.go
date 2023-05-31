package handler

import (
	"Mou1ght-Server/api/router"
	"Mou1ght-Server/internal/controller"
	"Mou1ght-Server/internal/model"
	"Mou1ght-Server/package/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func registerHandler(c *gin.Context) {
	username := c.Param("name")
	nickname := username
	password := c.Param("password")
	user := model.User{
		Name:     username,
		NickName: nickname,
	}
	if len(password) < 6 {
		router.Response(c, http.StatusUnprocessableEntity, 422, nil, "Your password shorter than 6 digits")
		return
	}

	if duplicate, _ := controller.CheckExistName(&user, username); duplicate {
		router.Response(c, http.StatusConflict, 409, nil, "Already exist user")
		return
	} else {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			router.Response(c, http.StatusInternalServerError, 500, nil, "System error")
			return
		}
		user.Password = string(hashedPassword)
		err = controller.RegisterUser(&user)
		if err != nil {
			router.Response(c, http.StatusInternalServerError, 500, nil, "User register failed with database error")
			return
		}
		token, err := util.ReleaseToken(&user)
		if err != nil {
			router.Response(c, http.StatusInternalServerError, 500, nil, "System error")
			log.Println("token generate error:", err)
			return
		}

		// All passed
		router.Success(c, gin.H{"token": token}, "Register successfully")
		return
	}
}
