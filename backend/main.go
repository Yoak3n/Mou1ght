package main

import (
	"Mou1ght-Server/api/router"
	"Mou1ght-Server/internal/database"
)

//func init() {
//
//}

func main() {
	defer database.GetConn().Close()
	router.RunSever()
}
