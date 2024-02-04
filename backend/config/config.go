package config

import (
	"Mou1ght-Server/package/logger"
	"github.com/spf13/viper"
)

type Configuration struct {
	SeverPort      int    `yaml:"sever_port"`
	DBName         string `yaml:"mysql_name"`
	DBPassword     string `yaml:"mysql_password"`
	DBAddr         string `yaml:"mysql_addr"`
	DBPort         int    `yaml:"mysql_port"`
	DatabaseName   string `yaml:"database_name"`
	DatabaseOption string `yaml:"database_option"`
	JwtKey         []byte `yaml:"jwt_key"`
}

var Conf *Configuration

// database options
const (
	SQLITE3    = "sqlite3"
	MYSQL      = "mysql"
	POSTGRESQL = "postgresql"
)

func init() {
	v := viper.New()
	Conf = new(Configuration)
	v.SetDefault("SEVER_PORT", 10420)
	v.SetDefault("DB_NAME", "root")
	v.SetDefault("DB_ADDR", "127.0.0.1")
	v.SetDefault("DB_PORT", 3306)
	v.SetDefault("DB_NAME", "mou1ght")
	v.SetDefault("DB_OPTION", SQLITE3)
	v.SetDefault("JWT_KEY", "mou1ght")
	ok := loadFromFile(v)
	if !ok {
		loadFromEnv(v)
		logger.Trace.Println("Load configuration from environment")
	}
	// check database option
	if Conf.DatabaseOption != SQLITE3 && Conf.DatabaseOption != MYSQL {
		logger.Error.Printf("Please choose one database option :[%s,%s,%s]\n", SQLITE3, MYSQL, POSTGRESQL)
	}
	logger.Info.Println(Conf)

}

func loadFromEnv(v *viper.Viper) {
	err := v.BindEnv("SEVER_PORT", "MYSQL_NAME", "MYSQL_PASSWORD", "MYSQL_PORT", "DB_NAME", "JWT_KEY")
	if err != nil {
		logger.Info.Println("GET ENVIRONMENT VARIABLE FAILED")
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
			logger.Warning.Println("Config file not exists")
		} else {
			logger.Error.Println("Read config file error")
		}
		return false
	}
	loadConfig(v)
	logger.Info.Println("Load configuration from file")
	return true
}

func loadConfig(v *viper.Viper) {
	Conf.SeverPort = v.GetInt("SEVER_PORT")
	Conf.DBName = v.GetString("DB_NAME")
	Conf.DBPassword = v.GetString("DB_PASSWORD")
	Conf.DBAddr = v.GetString("DB_ADDR")
	Conf.DBPort = v.GetInt("DB_PORT")
	Conf.DatabaseName = v.GetString("DB_NAME")
	Conf.DatabaseOption = v.GetString("DB_OPTION")
	Conf.JwtKey = []byte(v.GetString("JWT_KEY"))
}
