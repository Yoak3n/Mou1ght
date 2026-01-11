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
	e := r.Listen(":10420")
	if e != nil {
		panic(e)
	}
}

func main() {
	runApp()
}
