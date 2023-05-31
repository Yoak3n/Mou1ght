package database

import (
	"Mou1ght-Server/config"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func initSqlite() *gorm.DB {
	dsn := fmt.Sprintf("./%s.db", config.Conf.DatabaseName)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("database connected err:", err)
	}
	if err != nil {
		log.Panic("Create table failed:", err)
	}
	return db
}
