package interfaces

import "Mou1ght/internal/domain/model/table"

type TagRepository interface {
	CreateTag(tag *table.TagTable) error
	UpdateTag(tag *table.TagTable) error
	DeleteTagWithLink(tagID string) error
	CreateTagLink(tagLink *table.TagLinkTable) error
	DeleteTagLink(linkID string) error
	GetArticlesFromTagLink(link *table.TagLinkTable, desc bool) ([]table.ArticleTable, error)
	GetTagLinkByKeyword(keyword []string, typ string) (map[string]table.TagTable, []table.TagLinkTable, error)
}
