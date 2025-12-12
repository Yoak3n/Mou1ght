package controller

import (
	"Mou1ght/internal/domain/entity"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/instance"
	"errors"
)

func CreateArticle(req *request.CreateArticleRequest) error {
	aid := util.GenArticleID()
	record := &table.ArticleTable{
		PostBase: table.PostBase{
			ID:      aid,
			Content: req.Content,
		},
		Title: req.Title,
	}
	err := instance.UseDatabase().CreateArticle(record)
	if err != nil {
		return err
	}
	tagIDs := make([]string, len(req.Tags))
	for i, tag := range req.Tags {
		tagIDs[i] = tag.ID
	}
	err = CreateTagsLinkToArticle(tagIDs, aid)
	if err != nil {
		return err
	}
	categoryIDs := make([]string, len(req.Categories))
	for i, category := range req.Categories {
		categoryIDs[i] = category.ID
	}
	err = CreateCategoriesLinkToArticle(categoryIDs, aid)
	if err != nil {
		return err
	}
	return nil
}

func UpdateArticle(req *request.UpdateArticleRequest) error {
	record := &table.ArticleTable{
		PostBase: table.PostBase{
			ID:      req.ID,
			Content: req.Content,
		},
		Title:    req.Title,
		AuthorID: req.Author,
	}
	err := instance.UseDatabase().UpdateArticle(record)
	if err != nil {
		return err
	}
	categoryIDs := make(map[string]bool)
	for _, category := range req.Categories {
		categoryIDs[category.ID] = true
	}
	err = instance.UseDatabase().UpdateCategoryLinks(req.ID, categoryIDs)
	if err != nil {
		return err
	}
	tagsIDs := make(map[string]bool)
	for _, tag := range req.Tags {
		tagsIDs[tag.ID] = true
	}
	err = instance.UseDatabase().UpdateTargetLinks(req.ID, 1, tagsIDs)
	if err != nil {
		return err
	}

	return nil
}

func ViewArticle(id string) error {
	return instance.UseDatabase().AddViewCountArticle(id)
}

func LikeArticle(id string) error {
	return instance.UseDatabase().AddLikeCountArticle(id)
}

func GetArticleByID(id string) (*entity.ArticleEntity, error) {
	record, err := instance.UseDatabase().GetArticleByID(id)
	if err != nil {
		return nil, err
	}
	e := entity.NewArticleEntityFromTable(record, true)
	if e == nil {
		return nil, errors.New("article not exist")
	}
	return e, nil
}

func DeleteArticleByID(id string) error {
	err := instance.UseDatabase().DeleteArticleByID(id)
	if err != nil {
		return err
	}
	err = instance.UseDatabase().DeleteTagLinkFromTarget(id, 1)
	if err != nil {
		return err
	}
	err = instance.UseDatabase().DeleteArticleByID(id)
	if err != nil {
		return err
	}
	return nil
}
