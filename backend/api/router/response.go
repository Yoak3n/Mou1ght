package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(c *gin.Context, httpStatusCode int, code int, data gin.H, message string) {
	c.JSON(httpStatusCode, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func Success(c *gin.Context, data gin.H, message string) {
	Response(c, http.StatusOK, 200, data, message)
}
func Fail(c *gin.Context, message string, data gin.H) {
	Response(c, http.StatusBadRequest, 400, data, message)
}
