package main

import (
	"Mou1ght/consts"
	"Mou1ght/internal/config"
	"Mou1ght/internal/service/router"
	"Mou1ght/pkg/util"
)

func runApp() {
	config.GetConfig()
	util.CreateDirNotExists(consts.Upload)
	r := router.InitRouter()
	r.Listen(":10420")
}

func main() {
	runApp()
}
