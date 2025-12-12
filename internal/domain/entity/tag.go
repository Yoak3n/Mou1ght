package entity

import "Mou1ght/internal/domain/model/table"

type TagWithArticlesEntity struct {
	Tag      TagInformationEntity `json:"tag"`
	Articles []ArticleEntity      `json:"articles"`
}

type TagWithSharingEntity struct {
	Tag      TagInformationEntity `json:"tag"`
	Sharings []SharingEntity      `json:"sharings"`
}

type TagInformationEntity struct {
	ID    string `json:"id"`
	Label string `json:"label"`
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

func NewTagInformationEntityFromTable(tag *table.TagTable) TagInformationEntity {
	return TagInformationEntity{
		ID:    tag.ID,
		Label: tag.Label,
	}
}
