package database

import (
	"Mou1ght-Server/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// Created at 2023/4/10 14:56
// Created by Yoake

func InitMysql() *gorm.DB {
	conf := config.Conf
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.MysqlName, conf.MysqlPassword, conf.MysqlPort, conf.DatabaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("database connected err:", err)
	}
	if err != nil {
		log.Panic("Create table failed:", err)
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
