package handler

import (
	"Mou1ght-Server/internal/controller"
	"Mou1ght-Server/internal/model"
	"Mou1ght-Server/package/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func registerHandler(c *gin.Context) {
	username := c.Param("name")
	nickname := username
	password := c.Param("password")
	user := model.User{
		Name:     username,
		NikeName: nickname,
		Password: password,
	}
	if len(password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Your password shorter than 6 digits",
		})
		return
	}

	if duplicate := controller.CheckDuplicateName(&user, username); duplicate {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Already exist user",
		})
		return
	} else {
		err := controller.RegisterUser(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "User register failed with database error",
			})
			return
		}
		token, err := util.ReleaseToken(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "System error",
			})
			log.Println("token generate error:", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Register successfully",
			"data": gin.H{
				"token": token,
			},
		})
		return
	}
}
