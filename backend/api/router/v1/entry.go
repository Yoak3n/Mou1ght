package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterRouterGroup(r *gin.Engine) {
	version1 := r.Group("/v1")
	registerUserRouter(version1)
	registerArticleRouter(version1)
}
