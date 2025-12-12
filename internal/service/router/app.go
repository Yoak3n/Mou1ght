package router

import (
	"Mou1ght/internal/api/handler"
	"Mou1ght/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func setupUserRouter(r fiber.Router) {
	user := r.Group("/user")
	user.Post("/register", handler.Register)
	user.Post("/login", handler.Login)
	user.Get("/info", handler.Info).Use(middleware.Auth)
	user.Post("/logout", handler.Logout).Use(middleware.Auth)
}

func setupSettingRouter(r fiber.Router) {
	setting := r.Group("/setting")
	setting.Get("/all", handler.GetAllSetting)
	setting.Get("/blog", handler.GetBlogSetting)
	setting.Post("/blog", handler.UpdateBlogSetting)
}
