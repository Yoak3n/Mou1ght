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
