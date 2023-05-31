package database

import (
	"Mou1ght-Server/config"
	"Mou1ght-Server/internal/model"
	"log"
	"time"

	"database/sql"
	"gorm.io/gorm"
)

var Conn *sql.DB
var DB *gorm.DB

func init() {
	switch config.Conf.DatabaseOption {
	case "sqlite3":
		DB = InitSqlite()
		log.Println("Already connected Sqlite3")
	case "mysql":
		DB = InitMysql()
		log.Println("Already connected Mysql")
	}
	_ = DB.AutoMigrate(&model.User{}, &model.Paper{})
	Conn, _ = DB.DB()

	Conn.SetMaxOpenConns(100)
	Conn.SetMaxIdleConns(10)
	Conn.SetConnMaxLifetime(time.Hour)
}
