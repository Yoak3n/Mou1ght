package config

import (
	"Mou1ght/internal/domain/model/schema/console"
	"log"
	"sync"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type Configuration struct {
	Blog     console.BlogSetting `yaml:"blog"`
	Database DatabaseSetting     `yaml:"database"`
}

var config *Configuration
var once sync.Once

func (c *Configuration) ToMap() map[string]any {
	return map[string]any{
		"blog": c.Blog.ToMap(),
	}
}

func ReadConfig() *Configuration {
	err := viper.ReadInConfig()
	if err != nil {
		createDefaultConfInFile()
	}
	c := &Configuration{}
	err = viper.Unmarshal(c, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "yaml"
	})
	if err != nil {
		panic(err)
	}
	return c
}

func DefaultConfig() *Configuration {
	return &Configuration{
		Blog:     console.DefaultBlogSetting(),
		Database: DefaultDatabaseSetting(),
	}
}

func createDefaultConfInFile() {
	viper.Set("blog", console.DefaultBlogSetting())
	viper.Set("database", DefaultDatabaseSetting())
	err := viper.SafeWriteConfig()
	if err != nil {
		panic(err)
	}
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func GetConfig() *Configuration {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		config = ReadConfig()
	})
	return config
}

func UpdateConfig(c *Configuration) {
	config = c
	viper.Set("blog", c.Blog)
	viper.Set("database", c.Database)
	err := viper.WriteConfig()
	if err != nil {
		log.Printf("Failed to update config: %v", err)
	}
}

func UpdateBlogSetting(bs console.BlogSetting) {
	config.Blog = bs
	UpdateConfig(config)
}
