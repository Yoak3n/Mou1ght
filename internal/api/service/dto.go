package service

import (
	"Mou1ght/internal/domain/entity"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/interfaces"
	"strings"
)

type DTOService struct {
	ur      interfaces.UserRepository
	ar      interfaces.ArticleRepository
	tr      interfaces.TagRepository
	cr      interfaces.CategoryRepository
	counter interfaces.PostCounter
}

func NewDTOService(ur interfaces.UserRepository, ar interfaces.ArticleRepository, tr interfaces.TagRepository, cr interfaces.CategoryRepository, counter interfaces.PostCounter) *DTOService {
	return &DTOService{ur: ur, ar: ar, tr: tr, cr: cr, counter: counter}
}

func (s *DTOService) GetArticleEntityFromTable(article *table.ArticleTable, detail bool) *entity.ArticleEntity {
	user, err := s.ur.GetUser(article.AuthorID)
	if err != nil {
		return nil
	}
	viewDelta, likeDelta := s.counter.GetCounterDelta("article", article.ID)
	length := util.MeasureArticleLength(article.Content)
	content := ""
	if detail {
		content = article.Content
	} else {
		content = util.GenerateBriefFromMarkdown(article.Content)
	}
	e := &entity.ArticleEntity{
		ID:      article.ID,
		Title:   article.Title,
		Content: content,
		State: entity.PostState{
			View:   article.View + viewDelta,
			Like:   article.Like + likeDelta,
			Length: length,
			Status: entity.StatusIntToString(article.Status),
		},
		Categories: make([]entity.PostSign, 0),
		Tags:       make([]entity.PostSign, 0),
		Time: entity.PostTimeInfo{
			CreatedAt: article.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: article.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		Author: *entity.NewUserEntityFromTable(user, false),
	}
	tags, err := s.tr.QueryTagsByID(article.ID, table.ArticleTag)
	if err == nil {
		e.Tags = entity.NewTagsInformationEntityFromTable(tags)
	}
	categories, err := s.cr.QueryCategoriesByArticleID(article.ID)
	if err == nil {
		e.Categories = s.GetCategoriesInformationEntityFromTable(categories)
	}
	return e
}

func (s *DTOService) GetArticlesEntiesFromTable(list []*table.ArticleTable, detail bool) []*entity.ArticleEntity {
	es := make([]*entity.ArticleEntity, len(list))
	for i, article := range list {
		es[i] = s.GetArticleEntityFromTable(article, detail)
	}
	return es
}

func (s *DTOService) GetCategoryWithArticlesEntityFromTable(category *table.CategoryTable, articles []table.ArticleTable) *entity.CategoryWithArticlesEntity {
	as := make([]entity.ArticleEntity, len(articles))
	for i, article := range articles {
		as[i] = *s.GetArticleEntityFromTable(&article, false)
	}
	return &entity.CategoryWithArticlesEntity{
		Category: entity.NewCategoryInformationEntityFromTable(category),
		Articles: as,
	}
}

func (s *DTOService) GetCategoriesInformationEntityFromTable(categories []table.CategoryTable) []entity.PostSign {
	ps := make([]entity.PostSign, len(categories))
	for i, category := range categories {
		ps[i] = entity.NewCategoryInformationEntityFromTable(&category)
	}
	return ps
}

func (s *DTOService) GetCategoryInformationEntityFromTable(items []table.CategoryTable) []*entity.CategoryGroup {
	return s.GetCategoryGroupFromTables(items)
}

func (s *DTOService) GetCategoryGroupFromTables(items []table.CategoryTable) []*entity.CategoryGroup {
	nodeMap := make(map[string]*entity.CategoryGroup)
	for _, item := range items {
		node := &entity.CategoryGroup{
			PostSign: entity.NewCategoryInformationEntityFromTable(&item),
			Parent:   item.ParentID,
			Children: make([]*entity.CategoryGroup, 0),
		}
		nodeMap[item.ID] = node
	}
	var rootNodes = make([]*entity.CategoryGroup, 0)
	for _, item := range items {
		node := nodeMap[item.ID]
		if item.ParentID != "" {
			if parent, ok := nodeMap[item.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		} else {
			rootNodes = append(rootNodes, node)
		}
	}
	return rootNodes
}

func (s *DTOService) GetMessageEntityFromTable(msg *table.MessageTable) *entity.MessageEntity {
	length := util.MeasureArticleLength(msg.Content)
	viewDelta, likeDelta := s.counter.GetCounterDelta("message", msg.ID)
	return &entity.MessageEntity{
		ID:      msg.ID,
		Content: msg.Content,
		Position: entity.MessagePosition{
			X: msg.X,
			Y: msg.Y,
			Z: msg.Z,
		},
		State: entity.PostState{
			Like:   msg.Like + likeDelta,
			View:   msg.View + viewDelta,
			Length: length,
			Status: entity.StatusIntToString(msg.Status),
		},
		AuthorIP: msg.AuthorIP,
		Time: entity.PostTimeInfo{
			CreatedAt: msg.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: msg.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}
}

func (s *DTOService) GetMessagesEntityFromTables(msgs []*table.MessageTable) []*entity.MessageEntity {
	entities := make([]*entity.MessageEntity, 0, len(msgs))
	for _, msg := range msgs {
		entity := s.GetMessageEntityFromTable(msg)
		entities = append(entities, entity)
	}
	return entities
}

func (s *DTOService) GetSharingEntityFromTable(sharing *table.SharingTable) *entity.SharingEntity {
	user, err := s.ur.GetUser(sharing.AuthorID)
	if err != nil {
		return nil
	}
	viewDelta, likeDelta := s.counter.GetCounterDelta("sharing", sharing.ID)
	length := util.MeasureArticleLength(sharing.Content)
	e := &entity.SharingEntity{
		ID:      sharing.ID,
		Content: sharing.Content,
		Author:  *entity.NewUserEntityFromTable(user, false),
		State: entity.PostState{
			Like:   sharing.Like + likeDelta,
			View:   sharing.View + viewDelta,
			Length: length,
			Status: entity.StatusIntToString(sharing.Status),
		},
		Time: entity.PostTimeInfo{
			CreatedAt: sharing.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: sharing.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		Attachments: entity.NewAttachmentsEntityFromPaths(strings.Split(sharing.Attachment, ",")),
	}
	tags, err := s.tr.QueryTagsByID(sharing.ID, table.SharingTag)
	if err == nil {
		e.Tags = entity.NewTagsInformationEntityFromTable(tags)
	}
	return e
}

func (s *DTOService) GetSharingsEntityFromTables(sharings []*table.SharingTable) []entity.SharingEntity {
	entities := make([]entity.SharingEntity, 0, len(sharings))
	for _, sharing := range sharings {
		entity := s.GetSharingEntityFromTable(sharing)
		if entity != nil {
			entities = append(entities, *entity)
		}
	}
	return entities
}

func (s *DTOService) GetTagWithArticlesEntityFromTable(tag *table.TagTable, articles []table.ArticleTable) *entity.TagWithArticlesEntity {
	as := make([]entity.ArticleEntity, len(articles))
	for i, article := range articles {
		as[i] = *s.GetArticleEntityFromTable(&article, false)
	}
	return &entity.TagWithArticlesEntity{
		Tag:      entity.NewTagInformationEntityFromTable(tag),
		Articles: as,
	}
}

func (s *DTOService) GetTagWithSharingEntityFromTable(tag *table.TagTable, sharings []table.SharingTable) *entity.TagWithSharingEntity {
	ss := make([]entity.SharingEntity, len(sharings))
	for i, sharing := range sharings {
		ss[i] = *s.GetSharingEntityFromTable(&sharing)
	}
	return &entity.TagWithSharingEntity{
		Tag:      entity.NewTagInformationEntityFromTable(tag),
		Sharings: ss,
	}
}

func (s *DTOService) GetUserWithPostEntityFromTable(author *table.UserTable, sharings []table.SharingTable, articles []table.ArticleTable) *entity.UserWithPostEntity {
	sharing := make([]entity.SharingEntity, 0)
	for _, st := range sharings {
		sharing = append(sharing, *s.GetSharingEntityFromTable(&st))
	}
	article := make([]entity.ArticleEntity, 0)
	for _, a := range articles {
		article = append(article, *s.GetArticleEntityFromTable(&a, false))
	}

	return &entity.UserWithPostEntity{
		Author:   entity.NewUserEntityFromTable(author, false),
		Sharings: sharing,
		Articles: article,
	}
}
