package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Configuration struct {
	SeverPort      int
	MysqlName      string
	MysqlPassword  string
	MysqlPort      int
	DatabaseName   string
	DatabaseOption string
	JwtKey         []byte
}

var Conf *Configuration

// database options
const (
	SQLITE3 = "sqlite3"
	MYSQL   = "mysql"
)

func init() {
	v := viper.New()
	Conf = new(Configuration)
	v.SetDefault("SEVER_PORT", 10420)
	v.SetDefault("MYSQL_NAME", "root")
	v.SetDefault("MYSQL_PORT", 3306)
	v.SetDefault("DB_NAME", "mou1ght")
	v.SetDefault("DB_OPTION", SQLITE3)
	v.SetDefault("JWT_KEY", "mou1ght")
	ok := loadFromFile(v)
	if !ok {
		loadFromEnv(v)
		log.Println("Load configuration from environment")
	}
	log.Println("Load configuration from file")
	// check database option
	if Conf.DatabaseOption != SQLITE3 && Conf.DatabaseOption != MYSQL {
		log.Panic("Please choose one database option :", fmt.Sprintf("[%s,%s]", SQLITE3, MYSQL))
	}
	fmt.Println("\n", string(Conf.JwtKey), "\n")

}

func loadFromEnv(v *viper.Viper) {
	err := v.BindEnv("SEVER_PORT", "MYSQL_NAME", "MYSQL_PASSWORD", "MYSQL_PORT", "DB_NAME", "JWT_KEY")
	if err != nil {
		log.Println("GET ENVIRONMENT VARIABLE FAILED")
	}
	v.AutomaticEnv()
	v.AllowEmptyEnv(true)
	loadConfig(v)
}

func loadFromFile(v *viper.Viper) (readed bool) {
	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.SetConfigType("yml")
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not exists")
		} else {
			log.Println("Read config file error")
		}
		return false
	}
	loadConfig(v)
	return true
}

func loadConfig(v *viper.Viper) {
	Conf.SeverPort = v.GetInt("SEVER_PORT")
	Conf.MysqlName = v.GetString("MYSQL_NAME")
	Conf.MysqlPassword = v.GetString("MYSQL_PASSWORD")
	Conf.MysqlPort = v.GetInt("MYSQL_PORT")
	Conf.DatabaseName = v.GetString("DB_NAME")
	Conf.DatabaseOption = v.GetString("DB_OPTION")
	Conf.JwtKey = []byte(v.GetString("JWT_KEY"))
}
