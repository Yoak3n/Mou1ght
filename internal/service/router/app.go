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
	user.Use(middleware.Auth).Get("/info", handler.Info)
	user.Use(middleware.Auth).Post("/logout", handler.Logout)
}

func setupSettingRouter(r fiber.Router) {
	setting := r.Group("/setting")
	setting.Get("/all", handler.GetAllSetting)
	setting.Get("/blog", handler.GetBlogSetting)
	setting.Put("/blog", handler.UpdateBlogSetting)
}

func setupAttachmentRouter(r fiber.Router) {
	attachment := r.Group("/attachment")
	attachment.Use(middleware.Auth).Post("/upload", handler.UploadAttachment)
	attachment.Use(middleware.Auth).Get("/list", handler.GetAttachmentList)
}
