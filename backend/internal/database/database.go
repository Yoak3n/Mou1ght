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
var mdb *gorm.DB

func init() {
	switch config.Conf.DatabaseOption {
	case "sqlite3":
		mdb = initSqlite()
		logger.INFO("Already connected to Sqlite3")
	case "mysql":
		mdb = initMysql()
		logger.INFO("Already connected to Mysql")
	}
	_ = mdb.AutoMigrate(&model.User{}, &model.Article{})
	conn, _ = mdb.DB()

	conn.SetMaxOpenConns(100)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(time.Hour)
}

func GetDB() *gorm.DB {
	return mdb
}
func GetConn() *sql.DB {
	return conn
}
