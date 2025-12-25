package console

type NavBar struct {
	Links              []Link             `yaml:"links" json:"links"`
	WebsiteInformation WebsiteInformation `yaml:"website_information" json:"website_information"`
}

type Link struct {
	// 链接类型，home 表示首页，category 表示分类链接，external 表示外部链接，tag 表示标签链接
	Type string `json:"type" yaml:"type"`
	// 链接地址
	Href string `json:"href" yaml:"href"`
	// 链接显示文本
	Label string `json:"label" yaml:"label"`
}

type WebsiteInformation struct {
	Title    string   `json:"title" yaml:"title"`
	Icon     string   `json:"icon" yaml:"icon"`
	Logo     string   `json:"logo" yaml:"logo"`
	Keywords []string `json:"keywords" yaml:"keywords"`
}
