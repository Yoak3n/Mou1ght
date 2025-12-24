package controller

import (
	"Mou1ght/internal/domain/entity"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/repository/instance"
	"time"
)

func SingleListWithPost(req *request.PostListRequest, typ string) map[string]any {
	return queryPost(req, typ)
}

func queryPost(req *request.PostListRequest, typ string) map[string]any {
	ret := make(map[string]any)
	var startDate, endDate *time.Time
	if req.Filter.DateRange != nil {
		sd, _ := time.Parse("2006-01-02 15:04:05", req.Filter.DateRange.StartDate)
		ed, _ := time.Parse("2006-01-02 15:04:05", req.Filter.DateRange.EndDate)
		startDate, endDate = &sd, &ed
	}
	switch typ {
	case "article":
		articles, err := instance.UseDatabase().GetArticles(startDate, endDate)
		if err == nil {
			ret["articles"] = entity.NewArticleEntityFromTableList(articles, false)
		}
	case "sharing":
		sharings, err := instance.UseDatabase().GetSharings(startDate, endDate)
		if err == nil {
			ret["sharings"] = entity.NewSharingsEntityFromTables(sharings)
		}
	case "message":
		messages, err := instance.UseDatabase().GetMessages(startDate, endDate)
		if err == nil {
			ret["messages"] = entity.NewMessagesEntityFromTables(messages)
		}
	}
	return ret
}

func AllListWithPost(req *request.PostListRequest) map[string]any {
	ret := make(map[string]any)
	ret["articles"] = queryPost(req, "article")["articles"]
	ret["sharings"] = queryPost(req, "sharing")["sharings"]
	ret["messages"] = queryPost(req, "message")["messages"]
	return ret
}
