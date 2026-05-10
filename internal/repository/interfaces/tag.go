package interfaces

import "Mou1ght/internal/domain/model/table"

type TagRepository interface {
	CreateTag(tag *table.TagTable) error
	DeleteTag(tagID string) error
	GetAllTags() ([]table.TagTable, error)
	GetTagByID(id string) (*table.TagTable, error)
	GetTagsByID(ids []string) ([]table.TagTable, error)
	UpdateTag(tag *table.TagTable) error
	UpdateTargetLinks(targetID string, targetType table.TagType, ids map[string]bool) error
	DeleteTagWithLink(tagID string) error
	CreateTagLink(tagLink *table.TagLinkTable) error
	QueryTagsByLabel(label []string) ([]table.TagTable, error)
	QueryTagsByID(targetID string, targetType table.TagType) ([]table.TagTable, error)
	DeleteTagLink(linkID string) error
	DeleteTagLinkByTagID(tagID string) error
	DeleteTagLinkFromTarget(targetID string, targetType table.TagType) error
	GetTagLinkByKeyword(keyword []string, typ string) (map[string]table.TagTable, []table.TagLinkTable, error)
	GetArticlesFromTagLink(link *table.TagLinkTable, desc bool) ([]table.ArticleTable, error)
	GetSharingFromTagLink(link *table.TagLinkTable, desc bool) ([]table.SharingTable, error)
	CreateTagsLinkToArticle(tags []string, articleID string) error
}
