package controller

import (
	"Mou1ght/internal/config"
	"Mou1ght/internal/domain/model/schema/console"
)

func GetAllSetting() (map[string]any, error) {
	return config.GetConfig().ToMap(), nil
}

func GetBlogSetting() (console.BlogSetting, error) {
	return config.GetConfig().Blog, nil
}

func UpdateBlogSetting(setting console.BlogSetting) error {
	config.UpdateBlogSetting(setting)
	return nil
}
