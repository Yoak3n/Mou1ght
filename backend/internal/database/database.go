package database

import (
	"Mou1ght-Server/config"
	"Mou1ght-Server/internal/model"
	"Mou1ght-Server/package/logger"
	"time"

	"database/sql"
	"gorm.io/gorm"
)

var conn *sql.DB
var dbc *gorm.DB

func init() {
	switch config.Conf.DatabaseOption {
	case "sqlite3":
		dbc = initSqlite()
		logger.LogOut("Already connected Sqlite3")
	case "mysql":
		dbc = initMysql()
		logger.LogOut("Already connected Mysql")
	}
	_ = dbc.AutoMigrate(&model.User{}, &model.Article{})
	conn, _ = dbc.DB()

	conn.SetMaxOpenConns(100)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(time.Hour)
}

func GetDB() *gorm.DB {
	return dbc
}
func GetConn() *sql.DB {
	return conn
}
