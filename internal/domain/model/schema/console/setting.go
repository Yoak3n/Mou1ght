package console

type BlogSetting struct {
	NavBar
	BottomExtra
}

func DefaultBlogSetting() BlogSetting {
	return BlogSetting{
		NavBar: NavBar{
			SearchBox: "",
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
		"search_box": bs.SearchBox,
		"links":      bs.Links,
		"title":      bs.WebsiteInformation.Title,
		"icon":       bs.WebsiteInformation.Icon,
		"logo":       bs.WebsiteInformation.Logo,
		"keywords":   bs.WebsiteInformation.Keywords,
		"html":       bs.BottomExtra.HTML,
		"css":        bs.BottomExtra.CSS,
	}
}
