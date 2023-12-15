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
	// gin.SetMode(gin.ReleaseMode)
	R = gin.Default()
	R.SetTrustedProxies([]string{"http://127.0.0.1:7890"})
}

func RunSever() {
	addr := fmt.Sprintf(":%d", config.Conf.SeverPort)
	R.Use(middleware.CORSMiddleware())
	router.RegisterRouterGroup(R)
	_ = R.Run(addr)
}
