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
		logger.Info.Println("Already connected to Sqlite3")
	case "mysql":
		mdb = initMysql()
		logger.Info.Println("Already connected to MySQL")
	}
	migrateTables(&model.Article{}, &model.User{})
	conn, _ = mdb.DB()

	conn.SetMaxOpenConns(100)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(time.Hour)
}

func migrateTables(tables ...interface{}) {
	err := mdb.AutoMigrate(tables...)
	if err != nil {
		logger.Error.Println(err)
		logger.Error.Panic("migrate tables failed")
	}
}

func GetDB() *gorm.DB {
	return mdb
}
func GetConn() *sql.DB {
	return conn
}
