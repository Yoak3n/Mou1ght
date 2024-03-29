package router

import (
	"Mou1ght-Server/api/middleware"
	"Mou1ght-Server/api/response"
	"Mou1ght-Server/internal/controller"
	"Mou1ght-Server/internal/dto"
	"Mou1ght-Server/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func registerArticleRouter(g *gin.RouterGroup) {
	u := g.Group("/article")
	u.GET("/info/:id", articleInfo)
	u.GET("/list", articleList)
	u.POST("/add", middleware.AuthMiddleware(), articleAdd)
}

func articleList(c *gin.Context) {
	articles, err := controller.GetArticleList()
	if err != nil {
		response.Fail(c, err.Error(), nil)
	} else {
		response.Success(c, gin.H{"articles": dto.ToArticleList(articles)}, "Get article list successfully")
	}
}

func articleAdd(c *gin.Context) {
	article := new(dto.ArticlePostDTO)
	err := c.BindJSON(&article)
	if err != nil {
		response.Fail(c, "Invalid article data", nil)
	}
	err = controller.AddArticle(article)
	if err != nil {
		response.Fail(c, err.Error(), nil)
	} else {
		response.Success(c, nil, "add article successfully")
	}

}

func articleInfo(c *gin.Context) {
	id := c.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		response.Fail(c, "Invalid article id", nil)
		return
	}
	var article model.Article

	ok, _ := controller.CheckExistArticle(&article, uint(atoi))
	if ok {
		response.Success(c, gin.H{
			"article": dto.ToArticleDTO(&article),
		}, "Get article successfully")
	} else {
		response.Response(c, http.StatusNoContent, 404, nil, "Not found this article")
	}
}
