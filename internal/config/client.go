package config

// ClientSetting 描述与前台 Next.js 客户端相关的配置。
type ClientSetting struct {
	// RevalidateURL 是前台按需失效缓存的 webhook 地址，留空表示禁用。
	RevalidateURL string `yaml:"revalidate_url"`
	// RevalidateSecret 与前台环境变量 REVALIDATE_SECRET 保持一致。
	RevalidateSecret string `yaml:"revalidate_secret"`
}

func DefaultClientSetting() ClientSetting {
	return ClientSetting{}
}
