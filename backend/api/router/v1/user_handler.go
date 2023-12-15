package router

import (
	"Mou1ght-Server/api/middleware"
	"Mou1ght-Server/api/response"
	"Mou1ght-Server/internal/controller"
	"Mou1ght-Server/internal/database"
	"Mou1ght-Server/internal/dto"
	"Mou1ght-Server/internal/model"
	"Mou1ght-Server/package/logger"
	"Mou1ght-Server/package/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	u.POST("/logout", middleware.AuthMiddleware(), logoutHandler)
}

func loginHandler(c *gin.Context) {
	username := c.Param("name")
	password := c.Param("password")
	logger.INFO(username)
	var user model.User
	// check user exists or not
	db.Where("name = ?", username).First(&user)
	if user.ID == 0 {
		response.Response(c, http.StatusUnauthorized, 403, nil, "User doesn't exist")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "Your password shorter than 6 digits")
		return
	}
	// Judge password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		response.Response(c, http.StatusUnauthorized, 403, nil, "Incorrect password")
		return
	}
	token, err := util.ReleaseToken(&user)
	if err != nil {
		response.Fail(c, err.Error(), nil)
	}
	response.Success(c, gin.H{"token": token}, "Login successfully")
}
func userInfoHandler(c *gin.Context) {
	user, ok := c.Get("User")
	if ok {
		response.Success(c, gin.H{"user": dto.ToUserDTO(user.(*model.User))}, "")
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
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "Your password shorter than 6 digits")
		return
	}

	if duplicate, _ := controller.CheckExistName(&user, username); duplicate {
		response.Response(c, http.StatusConflict, 409, nil, "Already exist user")
		return
	} else {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			response.Response(c, http.StatusInternalServerError, 500, nil, "System error")
			return
		}
		user.Password = string(hashedPassword)
		err = controller.RegisterUser(&user)
		if err != nil {
			response.Response(c, http.StatusInternalServerError, 500, nil, "User register failed with database error")
			return
		}
		token, err := util.ReleaseToken(&user)
		if err != nil {
			response.Response(c, http.StatusInternalServerError, 500, nil, "System error")
			logger.Error.Printf("token generate error:%v\n", err)
			return
		}

		// All passed
		response.Success(c, gin.H{"token": token}, "Register successfully")
		return
	}
}

func logoutHandler(c *gin.Context) {

}
