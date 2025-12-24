package entity

import "Mou1ght/internal/domain/model/table"

type TagWithArticlesEntity struct {
	Tag      PostSign        `json:"tag"`
	Articles []ArticleEntity `json:"articles"`
}

type TagWithSharingEntity struct {
	Tag      PostSign        `json:"tag"`
	Sharings []SharingEntity `json:"sharings"`
}

func NewTagWithArticlesEntityFromTable(tag *table.TagTable, articles []table.ArticleTable) *TagWithArticlesEntity {
	s := make([]ArticleEntity, len(articles))
	for i, article := range articles {
		s[i] = *NewArticleEntityFromTable(&article, false)
	}
	return &TagWithArticlesEntity{
		Tag:      NewTagInformationEntityFromTable(tag),
		Articles: s,
	}
}

func NewTagWithSharingEntityFromTable(tag *table.TagTable, sharings []table.SharingTable) *TagWithSharingEntity {
	s := make([]SharingEntity, len(sharings))
	for i, sharing := range sharings {
		s[i] = *NewSharingEntityFromTable(&sharing)
	}
	return &TagWithSharingEntity{
		Tag:      NewTagInformationEntityFromTable(tag),
		Sharings: s,
	}
}

func NewTagInformationEntityFromTable(tag *table.TagTable) PostSign {
	return PostSign{
		ID:    tag.ID,
		Label: tag.Label,
	}
}

func NewTagsInformationEntityFromTable(tags []table.TagTable) []PostSign {
	s := make([]PostSign, len(tags))
	for i, tag := range tags {
		s[i] = NewTagInformationEntityFromTable(&tag)
	}
	return s
}
