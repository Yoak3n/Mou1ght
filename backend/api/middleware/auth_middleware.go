package middleware

import (
	"Mou1ght-Server/api/router"
	"Mou1ght-Server/internal/database"
	"Mou1ght-Server/internal/model"
	"Mou1ght-Server/package/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		fmt.Println(tokenString)
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {

			log.Println("Format of token is incorrect")
			c.Abort()
		}
		tokenString = tokenString[7:]

		token, claims, err := util.ParseToken(tokenString)
		if err != nil || !token.Valid {
			router.Response(c, http.StatusUnauthorized, 401, nil, "Unauthorized")
			log.Println("Token is invalid")
			c.Abort()
			return
		}
		userID := claims.UID
		if userID == 0 {
			router.Response(c, http.StatusUnauthorized, 401, nil, "Unauthorized")
			c.Abort()
			return
		}

		// Authorized
		var user model.User
		database.GetDB().First(&user, userID)
		c.Set("User", &user)
		c.Next()
	}
}
