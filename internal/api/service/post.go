package service

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/repository/interfaces"
	"errors"
	"time"
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
}

func NewPostService(articles interfaces.ArticleRepository, sharings interfaces.SharingRepository, messages interfaces.MessageRepository, posts interfaces.PostRepository) *PostService {
	return &PostService{articles: articles, sharings: sharings, messages: messages, posts: posts}
}

func (ps *PostService) SingleListWithPost(req *request.PostListRequest, typ string) *PostResult {
	return ps.queryPost(req, typ)
}

func (ps *PostService) queryPost(req *request.PostListRequest, typ string) *PostResult {
	ret := &PostResult{}
	var startDate, endDate *time.Time
	if req.Filter.DateRange != nil {
		sd, _ := time.Parse("2006-01-02 15:04:05", req.Filter.DateRange.StartDate)
		ed, _ := time.Parse("2006-01-02 15:04:05", req.Filter.DateRange.EndDate)
		startDate, endDate = &sd, &ed
	}
	switch typ {
	case "article":
		articles, err := ps.articles.GetArticles(startDate, endDate)
		if err == nil {
			ret.Articles = articles
		}
	case "sharing":
		sharings, err := ps.sharings.GetSharings(startDate, endDate)
		if err == nil {
			ret.Sharings = sharings
		}
	case "message":
		messages, err := ps.messages.GetMessages(startDate, endDate)
		if err == nil {
			ret.Messages = messages
		}
	}
	return ret
}

func (ps *PostService) AllListWithPost(req *request.PostListRequest) *PostResult {
	ret := &PostResult{}
	ret.Articles = ps.queryPost(req, "article").Articles
	ret.Sharings = ps.queryPost(req, "sharing").Sharings
	ret.Messages = ps.queryPost(req, "message").Messages
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
	// pending only for message
	// case "pending":
	// 	status = 3
	default:
		return errors.New("invalid status")
	}
	return ps.posts.UpdatePostStatus(req.PostType, req.ID, status)
}
