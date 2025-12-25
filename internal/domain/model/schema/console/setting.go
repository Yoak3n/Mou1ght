package console

type BlogSetting struct {
	NavBar      NavBar      `yaml:"nav_bar" json:"nav_bar"`
	BottomExtra BottomExtra `yaml:"bottom_extra" json:"bottom_extra"`
	Board       Board       `yaml:"board" json:"board"`
}

func DefaultBlogSetting() BlogSetting {
	return BlogSetting{
		NavBar: NavBar{
			Links: []Link{
				{
					Type:  "home",
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
		Board: Board{
			Question:     "",
			Answer:       "",
			NeedReviewed: false,
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
