package console

type NavBar struct {
	SearchBox          string             `yaml:"search_box"`
	Links              []Link             `yaml:"links"`
	WebsiteInformation WebsiteInformation `yaml:"website_information"`
}

type Link struct {
	Typ   string `json:"type" yaml:"type"`
	Href  string `json:"href" yaml:"href"`
	Label string `json:"label" yaml:"label"`
}

type WebsiteInformation struct {
	Title    string   `json:"title" yaml:"title"`
	Icon     string   `json:"icon" yaml:"icon"`
	Logo     string   `json:"logo" yaml:"logo"`
	Keywords []string `json:"keywords" yaml:"keywords"`
}
