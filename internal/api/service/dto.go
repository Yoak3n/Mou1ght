package service

import (
	"Mou1ght/internal/domain/entity"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/interfaces"
)

type DTOService struct {
	ur          interfaces.UserRepository
	ar          interfaces.ArticleRepository
	tr          interfaces.TagRepository
	cr          interfaces.CategoryRepository
	attachments interfaces.AttachmentRepository
	sal         interfaces.SharingAttachmentLinkRepository
	counter     interfaces.PostCounter
}

func NewDTOService(ur interfaces.UserRepository, ar interfaces.ArticleRepository, tr interfaces.TagRepository, cr interfaces.CategoryRepository, attachments interfaces.AttachmentRepository, sal interfaces.SharingAttachmentLinkRepository, counter interfaces.PostCounter) *DTOService {
	return &DTOService{ur: ur, ar: ar, tr: tr, cr: cr, attachments: attachments, sal: sal, counter: counter}
}

func missingUserEntity(id string) entity.UserEntity {
	return entity.UserEntity{
		ID:       id,
		UserName: "unknown",
		Avatar:   "",
		Role:     "unknown",
	}
}

// attachmentEntitiesByIDs 按给定顺序组装附件实体；查询出错时返回空切片
func (s *DTOService) attachmentEntitiesByIDs(ids []string) []entity.AttachmentEntity {
	entities := make([]entity.AttachmentEntity, 0)
	if len(ids) == 0 {
		return entities
	}
	records, err := s.attachments.GetAttachmentsByIDs(ids)
	if err != nil {
		return entities
	}
	recordsMap := make(map[string]table.AttachmentTable, len(records))
	for _, r := range records {
		recordsMap[r.ID] = r
	}
	for _, id := range ids {
		if r, ok := recordsMap[id]; ok {
			entities = append(entities, attachmentEntityFromTable(&r))
		}
	}
	return entities
}

func (s *DTOService) GetArticleEntityFromTable(article *table.ArticleTable, detail bool) *entity.ArticleEntity {
	if article == nil {
		return nil
	}
	user, err := s.ur.GetUser(article.AuthorID)
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
	}
	if err != nil || user == nil {
		e.Author = missingUserEntity(article.AuthorID)
	} else {
		e.Author = *entity.NewUserEntityFromTable(user, false)
	}
	tags, err := s.tr.QueryTagsByID(article.ID, table.ArticleTag)
	if err == nil {
		e.Tags = entity.NewTagsInformationEntityFromTable(tags)
	}
	categories, err := s.cr.QueryCategoriesByArticleID(article.ID)
	if err == nil {
		e.Categories = s.GetCategoriesInformationEntityFromTable(categories)
	}
	if detail && s.sal != nil {
		ids, err := s.sal.GetAttachmentIDsByArticleID(article.ID)
		if err == nil {
			e.Attachments = s.attachmentEntitiesByIDs(ids)
		}
	}
	return e
}

func (s *DTOService) GetArticlesEntiesFromTable(list []*table.ArticleTable, detail bool) []*entity.ArticleEntity {
	es := make([]*entity.ArticleEntity, 0, len(list))
	for _, article := range list {
		e := s.GetArticleEntityFromTable(article, detail)
		if e != nil {
			es = append(es, e)
		}
	}
	return es
}

func (s *DTOService) GetCategoryWithArticlesEntityFromTable(category *table.CategoryTable, articles []table.ArticleTable) *entity.CategoryWithArticlesEntity {
	as := make([]entity.ArticleEntity, 0, len(articles))
	for i := range articles {
		e := s.GetArticleEntityFromTable(&articles[i], false)
		if e != nil {
			as = append(as, *e)
		}
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
	// 统计失败时计数留空，不阻塞分类树返回
	counts, err := s.cr.CountArticlesGroupByCategory()
	if err != nil {
		counts = map[string]int64{}
	}
	nodeMap := make(map[string]*entity.CategoryGroup)
	for _, item := range items {
		node := &entity.CategoryGroup{
			PostSign: entity.NewCategoryInformationEntityFromTable(&item),
			Parent:   item.ParentID,
			Children: make([]*entity.CategoryGroup, 0),
		}
		node.Count = counts[item.ID]
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
	// TotalCount 自底向上聚合子分类的文章数
	var fillTotal func(n *entity.CategoryGroup) int64
	fillTotal = func(n *entity.CategoryGroup) int64 {
		total := n.Count
		for _, child := range n.Children {
			total += fillTotal(child)
		}
		n.TotalCount = total
		return total
	}
	for _, root := range rootNodes {
		fillTotal(root)
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
	if sharing == nil {
		return nil
	}
	user, err := s.ur.GetUser(sharing.AuthorID)
	viewDelta, likeDelta := s.counter.GetCounterDelta("sharing", sharing.ID)
	length := util.MeasureArticleLength(sharing.Content)
	var attachments []entity.AttachmentEntity
	ids, err := s.sal.GetAttachmentIDsBySharingID(sharing.ID)
	if err == nil {
		attachments = s.attachmentEntitiesByIDs(ids)
	}
	e := &entity.SharingEntity{
		ID:      sharing.ID,
		Content: sharing.Content,
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
		Attachments: attachments,
	}
	if err != nil || user == nil {
		e.Author = missingUserEntity(sharing.AuthorID)
	} else {
		e.Author = *entity.NewUserEntityFromTable(user, false)
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
	as := make([]entity.ArticleEntity, 0, len(articles))
	for i := range articles {
		e := s.GetArticleEntityFromTable(&articles[i], false)
		if e != nil {
			as = append(as, *e)
		}
	}
	return &entity.TagWithArticlesEntity{
		Tag:      entity.NewTagInformationEntityFromTable(tag),
		Articles: as,
	}
}

func (s *DTOService) GetTagWithSharingEntityFromTable(tag *table.TagTable, sharings []table.SharingTable) *entity.TagWithSharingEntity {
	ss := make([]entity.SharingEntity, 0, len(sharings))
	for i := range sharings {
		e := s.GetSharingEntityFromTable(&sharings[i])
		if e != nil {
			ss = append(ss, *e)
		}
	}
	return &entity.TagWithSharingEntity{
		Tag:      entity.NewTagInformationEntityFromTable(tag),
		Sharings: ss,
	}
}

func (s *DTOService) GetUserWithPostEntityFromTable(author *table.UserTable, sharings []table.SharingTable, articles []table.ArticleTable) *entity.UserWithPostEntity {
	sharing := make([]entity.SharingEntity, 0)
	for i := range sharings {
		e := s.GetSharingEntityFromTable(&sharings[i])
		if e != nil {
			sharing = append(sharing, *e)
		}
	}
	article := make([]entity.ArticleEntity, 0)
	for i := range articles {
		e := s.GetArticleEntityFromTable(&articles[i], false)
		if e != nil {
			article = append(article, *e)
		}
	}

	return &entity.UserWithPostEntity{
		Author:   entity.NewUserEntityFromTable(author, false),
		Sharings: sharing,
		Articles: article,
	}
}
