package service

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/interfaces"
)

type ArticleService struct {
	articles      interfaces.ArticleRepository
	users         interfaces.UserRepository
	categories    interfaces.CategoryRepository
	categoryLinks interfaces.CategoryLinkRepository
	tags          interfaces.TagRepository
}

func NewArticleService(articles interfaces.ArticleRepository, categories interfaces.CategoryRepository, categoryLinks interfaces.CategoryLinkRepository, tags interfaces.TagRepository) *ArticleService {
	return &ArticleService{articles: articles, categories: categories, categoryLinks: categoryLinks, tags: tags}
}

func (s *ArticleService) CreateArticle(req *request.CreateArticleRequest) error {
	aid := util.GenArticleID()
	record := &table.ArticleTable{
		PostBase: table.PostBase{
			ID:      aid,
			Content: req.Content,
		},
		Title: req.Title,
	}
	err := s.articles.CreateArticle(record)
	if err != nil {
		return err
	}
	tagIDs := make([]string, len(req.Tags))
	for i, tag := range req.Tags {
		tagIDs[i] = tag.ID
	}
	err = s.tags.CreateTagsLinkToArticle(tagIDs, aid)
	if err != nil {
		return err
	}
	categoryIDs := make([]string, len(req.Categories))
	for i, category := range req.Categories {
		categoryIDs[i] = category.ID
	}
	err = s.categoryLinks.CreateCategoriesLinkToArticle(categoryIDs, aid)
	if err != nil {
		return err
	}
	return nil
}

func (s *ArticleService) UpdateArticle(req *request.UpdateArticleRequest) error {
	record := &table.ArticleTable{
		PostBase: table.PostBase{
			ID:      req.ID,
			Content: req.Content,
		},
		Title:    req.Title,
		AuthorID: req.Author,
	}
	err := s.articles.UpdateArticle(record)
	if err != nil {
		return err
	}
	categoryIDs := make(map[string]bool)
	for _, category := range req.Categories {
		categoryIDs[category.ID] = true
	}
	err = s.categoryLinks.UpdateCategoryLinks(req.ID, categoryIDs)
	if err != nil {
		return err
	}
	tagsIDs := make(map[string]bool)
	for _, tag := range req.Tags {
		tagsIDs[tag.ID] = true
	}
	err = s.tags.UpdateTargetLinks(req.ID, 1, tagsIDs)
	if err != nil {
		return err
	}

	return nil
}

func (s *ArticleService) ViewArticle(id string) error {
	return s.articles.AddViewCountArticle(id)
}

func (s *ArticleService) LikeArticle(id string) error {
	return s.articles.AddLikeCountArticle(id)
}

func (s *ArticleService) GetArticleByID(id string) (*table.ArticleTable, error) {
	record, err := s.articles.GetArticleByID(id)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (s *ArticleService) DeleteArticleByID(id string) error {
	err := s.articles.DeleteArticleByID(id)
	if err != nil {
		return err
	}
	err = s.tags.DeleteTagLinkFromTarget(id, 1)
	if err != nil {
		return err
	}
	err = s.articles.DeleteArticleByID(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ArticleService) GetArticlesByAuthorID(authorID string, descend bool) ([]table.ArticleTable, error) {
	return s.articles.GetArticlesByAuthorID(authorID, descend)
}
