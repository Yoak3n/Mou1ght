package main

import (
	"Mou1ght/internal/config"
	"Mou1ght/internal/service/router"
)

func runApp() {
	config.GetConfig()
	r := router.InitRouter()
	r.Listen(":10420")
}

func main() {
	runApp()
}
