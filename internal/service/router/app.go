package router

import (
	"Mou1ght/internal/api/handler"
	"Mou1ght/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func setupUserRouter(r fiber.Router, userHandler *handler.UserHandler) {
	user := r.Group("/user")
	user.Post("/register", userHandler.Register)
	user.Get("/register/status", userHandler.RegisterStatus)
	user.Post("/login", userHandler.Login)
	user.Use(middleware.Auth).Get("/info", userHandler.Info)
	user.Use(middleware.Auth).Post("/logout", userHandler.Logout)
	user.Use(middleware.Auth).Put("/profile", userHandler.UpdateProfile)
	user.Use(middleware.Auth).Put("/password", userHandler.ChangePassword)
}

func setupSettingRouter(r fiber.Router) {
	setting := r.Group("/setting")
	setting.Get("/all", middleware.Auth, handler.GetAllSetting)
	setting.Get("/blog/public", handler.GetPublicBlogSetting)
	setting.Get("/blog", middleware.Auth, handler.GetBlogSetting)
	setting.Put("/blog", middleware.Auth, handler.UpdateBlogSetting)
}

func setupAttachmentRouter(r fiber.Router, attachmentHandler *handler.AttachmentHandler) {
	attachment := r.Group("/attachment")
	attachment.Post("/upload", middleware.Auth, attachmentHandler.UploadAttachment)
	attachment.Get("/list", middleware.Auth, attachmentHandler.GetAttachmentList)
	attachment.Delete("/delete/:id", middleware.Auth, attachmentHandler.DeleteAttachment)
}
