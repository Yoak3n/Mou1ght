package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(c *gin.Context, httpStatusCode int, code int, data gin.H, message string) {
	type Res struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}
	c.JSON(httpStatusCode, Res{Code: code, Message: message, Data: data})
}

func Success(c *gin.Context, data gin.H, message string) {
	Response(c, http.StatusOK, 200, data, message)
}
func Fail(c *gin.Context, message string, data gin.H) {
	Response(c, http.StatusBadRequest, 400, data, message)
	c.Abort()
}
