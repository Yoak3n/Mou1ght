package console

type BlogSetting struct {
	NavBar      NavBar      `yaml:"nav_bar" json:"nav_bar"`
	BottomExtra BottomExtra `yaml:"bottom_extra" json:"bottom_extra"`
}

func DefaultBlogSetting() BlogSetting {
	return BlogSetting{
		NavBar: NavBar{
			Links: []Link{
				{
					Typ:   "Home",
					Href:  "/",
					Label: "Home",
				},
			},
			WebsiteInformation: WebsiteInformation{
				Title:    "Mou1ght",
				Icon:     "",
				Logo:     "",
				Keywords: []string{"Mou1ght", "Blog"},
			},
		},
		BottomExtra: BottomExtra{
			HTML: "",
			CSS:  "",
		},
	}
}

func (bs *BlogSetting) ToMap() map[string]any {
	return map[string]any{
		"links":    bs.NavBar.Links,
		"title":    bs.NavBar.WebsiteInformation.Title,
		"icon":     bs.NavBar.WebsiteInformation.Icon,
		"logo":     bs.NavBar.WebsiteInformation.Logo,
		"keywords": bs.NavBar.WebsiteInformation.Keywords,
		"html":     bs.BottomExtra.HTML,
		"css":      bs.BottomExtra.CSS,
	}
}
