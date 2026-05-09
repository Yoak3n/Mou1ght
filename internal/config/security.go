package config

type SecuritySetting struct {
	JwtKey        string `yaml:"jwt_key"`
	VisitorJwtKey string `yaml:"visitor_jwt_key"`
}

func DefaultSecuritySetting() SecuritySetting {
	return SecuritySetting{
		JwtKey:        "",
		VisitorJwtKey: "",
	}
}
