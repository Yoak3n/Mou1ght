package service

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/notify"
	"Mou1ght/internal/repository/interfaces"
	"errors"
)

type PostService struct {
	articles interfaces.ArticleRepository
	sharings interfaces.SharingRepository
	messages interfaces.MessageRepository
	posts    interfaces.PostRepository
}

type PostResult struct {
	Articles []*table.ArticleTable
	Sharings []*table.SharingTable
	Messages []*table.MessageTable
	Total    int64
	Page     int
	PageSize int
}

func NewPostService(articles interfaces.ArticleRepository, sharings interfaces.SharingRepository, messages interfaces.MessageRepository, posts interfaces.PostRepository) *PostService {
	return &PostService{articles: articles, sharings: sharings, messages: messages, posts: posts}
}

func (ps *PostService) SingleListWithPost(req *request.PostListRequest, typ string) *PostResult {
	return ps.queryPost(req, typ)
}

func (ps *PostService) queryPost(req *request.PostListRequest, typ string) *PostResult {
	opts := request.NewListOptions(req.Filter)
	ret := &PostResult{
		Page:     opts.Page,
		PageSize: opts.PageSize,
	}
	switch typ {
	case "article":
		articles, total, err := ps.articles.GetArticles(opts)
		if err == nil {
			ret.Articles = articles
			ret.Total = total
		}
	case "sharing":
		sharings, total, err := ps.sharings.GetSharings(opts)
		if err == nil {
			ret.Sharings = sharings
			ret.Total = total
		}
	case "message":
		messages, total, err := ps.messages.GetMessages(opts)
		if err == nil {
			ret.Messages = messages
			ret.Total = total
		}
	}
	return ret
}

func (ps *PostService) AllListWithPost(req *request.PostListRequest) *PostResult {
	opts := request.NewListOptions(req.Filter)
	ret := &PostResult{
		Page:     opts.Page,
		PageSize: opts.PageSize,
	}

	articles, articleTotal, err := ps.articles.GetArticles(opts)
	if err == nil {
		ret.Articles = articles
		ret.Total += articleTotal
	}
	sharings, sharingTotal, err := ps.sharings.GetSharings(opts)
	if err == nil {
		ret.Sharings = sharings
		ret.Total += sharingTotal
	}
	messages, messageTotal, err := ps.messages.GetMessages(opts)
	if err == nil {
		ret.Messages = messages
		ret.Total += messageTotal
	}
	return ret
}

func (ps *PostService) UpdatePostStatus(req *request.UpdatePostStatusRequest) error {
	status := int8(0)
	switch req.Status {
	case "draft":
		status = 0
	case "publish":
		status = 1
	case "archive":
		status = 2
	case "pending":
		if req.PostType != "message" {
			return errors.New("pending is only allowed for message")
		}
		status = 3
	default:
		return errors.New("invalid status")
	}
	if err := ps.posts.UpdatePostStatus(req.PostType, req.ID, status); err != nil {
		return err
	}
	notify.RevalidateClient()
	return nil
}