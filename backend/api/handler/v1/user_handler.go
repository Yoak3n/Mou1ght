package handler

import (
	"Mou1ght-Server/api/middleware"
	"Mou1ght-Server/api/router"
	"Mou1ght-Server/internal/controller"
	"Mou1ght-Server/internal/database"
	"Mou1ght-Server/internal/dto"
	"Mou1ght-Server/internal/model"
	"Mou1ght-Server/package/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var db *gorm.DB

func init() {
	db = database.GetDB()
}

func registerUserRouter(g *gin.RouterGroup) {
	u := g.Group("/user")
	u.POST("/login/:name/:password", loginHandler)
	u.POST("/register/:name/:password", registerHandler)
	u.GET("/info", middleware.AuthMiddleware(), userInfoHandler)
}

func loginHandler(c *gin.Context) {
	username := c.Param("name")
	password := c.Param("password")
	var user model.User
	// 查询是否存在
	db.Where("name = ?", username).First(&user)
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
