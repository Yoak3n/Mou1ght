package controller

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/repository/instance"
	"time"
)

func AllListWithPost(req *request.PostListRequest) map[string]any {
	ret := make(map[string]any)
	if req.Filter.DateRange != nil {
		articles, err := instance.UseDatabase().GetArticles(nil, nil)
		if err == nil {
			ret["articles"] = articles
		}
		sharings, err := instance.UseDatabase().GetSharings(nil, nil)
		if err == nil {
			ret["sharings"] = sharings
		}

	} else {
		var startDate, endDate time.Time
		startDate, _ = time.Parse("2006-01-02 15:04:05", req.Filter.DateRange.StartDate)
		endDate, _ = time.Parse("2006-01-02 15:04:05", req.Filter.DateRange.EndDate)
		articles, err := instance.UseDatabase().GetArticles(&startDate, &endDate)
		if err == nil {
			ret["articles"] = articles
		}
		sharings, err := instance.UseDatabase().GetSharings(&startDate, &endDate)
		if err == nil {
			ret["sharings"] = sharings
		}
	}

	return ret
}
