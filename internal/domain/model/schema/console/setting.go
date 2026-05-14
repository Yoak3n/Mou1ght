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

type PublicBoard struct {
	Question     string `json:"question"`
	NeedReviewed bool   `json:"need_reviewed"`
}

type PublicBlogSetting struct {
	NavBar      NavBar       `json:"nav_bar"`
	BottomExtra BottomExtra  `json:"bottom_extra"`
	Board       PublicBoard  `json:"board"`
}

func (bs *BlogSetting) ToPublic() PublicBlogSetting {
	return PublicBlogSetting{
		NavBar:      bs.NavBar,
		BottomExtra: bs.BottomExtra,
		Board: PublicBoard{
			Question:     bs.Board.Question,
			NeedReviewed: bs.Board.NeedReviewed,
		},
	}
}
