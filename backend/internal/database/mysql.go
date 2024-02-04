package database

import (
	"Mou1ght-Server/config"
	"Mou1ght-Server/package/logger"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Created at 2023/4/10 14:56
// Created by Yoake

func initMysql() *gorm.DB {
	conf := config.Conf
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.DBName, conf.DBPassword, conf.DBPort, conf.DatabaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error.Printf("database connected err:%v\n", err)
	}
	if err != nil {
		logger.Error.Printf("database connected err:%v\n", err)
	}
	return db
}

//func UserRegister(session string, uid string) {
//	DB.Create(&model.User{UID: uid, Session: session})
//}
//
//func ReadUser(session string) model.User {
//	DB.First(user, "session = ?", session)
//	return user
//}

//func CloseConnect() {
//	conn.Close()
//}
