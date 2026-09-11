package config

type DatabaseSetting struct {
	DSN  string `yaml:"dsn"`
	Type string `yaml:"type"`
}

func DefaultDatabaseSetting() DatabaseSetting {
	return DatabaseSetting{
		// 默认 SQLite 库文件放在 data 目录下，便于容器内持久化（/app/data 卷）
		DSN:  "data/Mou1ght",
		Type: "sqlite",
	}
}
