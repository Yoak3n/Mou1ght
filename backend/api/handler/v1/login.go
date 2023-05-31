package handler

import (
	"Mou1ght-Server/api/dto"
	"Mou1ght-Server/api/router"
	"Mou1ght-Server/internal/database"
	"Mou1ght-Server/internal/model"
	"Mou1ght-Server/package/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func loginHandler(c *gin.Context) {
	username := c.Param("name")
	password := c.Param("password")
	var user model.User
	// 查询是否存在

	database.DB.Where("name = ?", username).First(&user)
	if user.ID == 0 {
		router.Response(c, http.StatusNotAcceptable, 404, nil, "User doesn't exist")
		return
	}
	if len(password) < 6 {
		router.Response(c, http.StatusUnprocessableEntity, 422, nil, "Your password shorter than 6 digits")
		return
	}
	// Judge password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		router.Fail(c, "Incorrect password", nil)
		return
	}
	token, err := util.ReleaseToken(&user)
	router.Success(c, gin.H{"token": token}, "Login in successfully")
}
func userInfoHandler(c *gin.Context) {
	user, ok := c.Get("User")
	if ok {
		router.Success(c, gin.H{"user": dto.ToUserDTO(user.(*model.User))}, "")
	}
}
