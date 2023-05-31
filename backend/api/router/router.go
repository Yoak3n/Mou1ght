package router

import (
	"Mou1ght-Server/api/handler/v1"
	"Mou1ght-Server/config"
	"fmt"
	"github.com/gin-gonic/gin"
)

var R *gin.Engine

func init() {
	R = gin.Default()
}

func RunSever() {
	addr := fmt.Sprintf(":%d", config.Conf.SeverPort)
	handler.RegisterRouterGroup(R)
	_ = R.Run(addr)
}
