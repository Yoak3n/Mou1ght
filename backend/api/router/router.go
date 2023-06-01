package router

import (
	"Mou1ght-Server/api/middleware"
	"Mou1ght-Server/api/router/v1"
	_ "Mou1ght-Server/api/router/v1"
	"Mou1ght-Server/config"
	"fmt"
	"github.com/gin-gonic/gin"
)

var R *gin.Engine

func init() {
	//gin.SetMode(gin.ReleaseMode)
	R = gin.Default()
}

func RunSever() {
	addr := fmt.Sprintf(":%d", config.Conf.SeverPort)
	router.RegisterRouterGroup(R)
	R.Use(middleware.CORSMiddleware())
	_ = R.Run(addr)
}
