package database

import (
	"Mou1ght-Server/config"
	"Mou1ght-Server/package/logger"
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func initSqlite() *gorm.DB {

	dsn := fmt.Sprintf("./data/db/%s.db", config.Conf.DatabaseName)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error.Println(fmt.Sprintf("database connected err:%v", err))
	}
	if err != nil {
		logger.Error.Println(fmt.Sprintf("database connected err:%v", err))
	}
	return db
}
