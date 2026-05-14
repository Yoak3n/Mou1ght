package router

import (
	"Mou1ght/internal/api/handler"
	"Mou1ght/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func setupUserRouter(r fiber.Router, userHandler *handler.UserHandler) {
	user := r.Group("/user")
	user.Post("/register", userHandler.Register)
	user.Post("/login", userHandler.Login)
	user.Use(middleware.Auth).Get("/info", userHandler.Info)
	user.Use(middleware.Auth).Post("/logout", userHandler.Logout)
}

func setupSettingRouter(r fiber.Router) {
	setting := r.Group("/setting")
	setting.Use(middleware.Auth).Get("/all", handler.GetAllSetting)
	setting.Get("/blog/public", handler.GetPublicBlogSetting)
	setting.Use(middleware.Auth).Get("/blog", handler.GetBlogSetting)
	setting.Use(middleware.Auth).Put("/blog", handler.UpdateBlogSetting)
}

func setupAttachmentRouter(r fiber.Router, attachmentHandler *handler.AttachmentHandler) {
	attachment := r.Group("/attachment")
	attachment.Use(middleware.Auth).Post("/upload", attachmentHandler.UploadAttachment)
	attachment.Use(middleware.Auth).Get("/list", attachmentHandler.GetAttachmentList)
}
