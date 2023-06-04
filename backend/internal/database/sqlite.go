package database

import (
	"Mou1ght-Server/config"
	"Mou1ght-Server/package/logger"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initSqlite() *gorm.DB {
	dsn := fmt.Sprintf("./%s.db", config.Conf.DatabaseName)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.ERROR(fmt.Sprintf("database connected err:%v", err))
	}
	if err != nil {
		logger.ERROR(fmt.Sprintf("database connected err:%v", err))
	}
	return db
}
