package config

type DatabaseSetting struct {
	DSN  string `yaml:"dsn"`
	Type string `yaml:"type"`
}

func DefaultDatabaseSetting() DatabaseSetting {
	return DatabaseSetting{
		DSN:  "Mou1ght",
		Type: "sqlite",
	}
}
