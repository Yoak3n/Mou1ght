package router

import (
	"Mou1ght-Server/api/response"
	"Mou1ght-Server/internal/controller"
	"Mou1ght-Server/internal/dto"
	"Mou1ght-Server/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func registerArticleRouter(g *gin.RouterGroup) {
	u := g.Group("/article")
	//u.POST("/login/:name/:password", loginHandler)
	//u.POST("/register/:name/:password", registerHandler)
	u.GET("/info/:id", articleInfo)
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
