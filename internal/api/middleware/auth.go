package middleware

import (
	"Mou1ght/internal/pkg/util"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Auth(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	tokenHeader, ok := headers["Authorization"]
	if !ok || len(tokenHeader) < 1 || !strings.HasPrefix(tokenHeader[0], "Bearer ") {
		return util.ErrorResponse(c, 401, "Unauthorized")
	}
	tokenString := tokenHeader[0][7:]
	token, claims, err := util.ParseToken(tokenString)
	if err != nil || !token.Valid {
		return util.ErrorResponse(c, 401, "Unauthorized")
	}
	userId := claims.MapClaims["uid"].(string)
	if userId == "" {
		return util.ErrorResponse(c, 401, "Unauthorized")
	}
	c.Locals("uid", userId)
	c.Locals("token", tokenString)
	return c.Next()
}
