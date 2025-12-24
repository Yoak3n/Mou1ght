package console

type Board struct {
	Question     string `yaml:"question" json:"question"`
	Answer       string `yaml:"answer" json:"answer"`
	NeedReviewed bool   `yaml:"need_reviewed" json:"need_reviewed"`
}
