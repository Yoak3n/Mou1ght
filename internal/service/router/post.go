package router

import (
	"Mou1ght/internal/api/handler"
	"Mou1ght/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func setupArticleRouter(r fiber.Router) {
	article := r.Group("/article")
	article.Post("/create", handler.CreateArticle).Use(middleware.Auth)
	article.Delete("/delete/:id", handler.DeleteArticle).Use(middleware.Auth)
	article.Put("/update", handler.UpdateArticle).Use(middleware.Auth)
	article.Get("/list", handler.ListPost)
	article.Get("/detail/:id", handler.DetailArticle)
	article.Post("/view/:id", handler.ViewPost)
	article.Post("/like/:id", handler.LikePost)
}

func setupSharingRouter(r fiber.Router) {
	sharing := r.Group("/sharing")
	sharing.Post("/create", handler.CreateSharing).Use(middleware.Auth)
	sharing.Delete("/delete/:id", handler.DeleteSharing).Use(middleware.Auth)
	sharing.Put("/update", handler.UpdateSharing).Use(middleware.Auth)
	sharing.Get("/list", handler.ListPost)
	sharing.Get("/detail/:id", handler.DetailSharing)
	sharing.Post("/view/:id", handler.ViewPost)
	sharing.Post("/like/:id", handler.LikePost)
}

func setupMessageRouter(r fiber.Router) {
	message := r.Group("/message")
	message.Post("/create", handler.CreateMessage).Use(middleware.Auth)
	message.Delete("/delete/:id", handler.DeleteMessage).Use(middleware.Auth)
	message.Put("/update", handler.UpdateMessage).Use(middleware.Auth)
	message.Get("/list", handler.ListMessage)
	message.Get("/detail/:id", handler.DetailMessage)
	message.Post("/view/:id", handler.ViewMessage)
	message.Post("/like/:id", handler.LikeMessage)
}
