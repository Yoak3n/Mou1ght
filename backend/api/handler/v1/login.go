package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func loginHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"cookie": "123465",
	})
}
func userInfoHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name": "Yoake",
	})
}
