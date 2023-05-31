package handler

import (
	"Mou1ght-Server/api/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRouterGroup(r *gin.Engine) {
	Version1 := r.Group("/v1")
	Version1.POST("/login/:name/:password", loginHandler)
	Version1.POST("/register/:name/:password", registerHandler)
	Version1.GET("/info", middleware.AuthMiddleware(), userInfoHandler)
}
