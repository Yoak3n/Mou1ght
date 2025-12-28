package router

import (
	"Mou1ght/internal/api/handler"
	"Mou1ght/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func setupArticleRouter(r fiber.Router) {
	article := r.Group("/article").Name("article.")
	article.Post("/create", middleware.Auth, handler.CreateArticle)
	article.Delete("/delete/:id", middleware.Auth, handler.DeleteArticle)
	article.Put("/update", middleware.Auth, handler.UpdateArticle)
	article.Post("/list", handler.ListPost).Name("list")
	article.Get("/detail/:id", handler.DetailArticle).Name("detail")
	article.Post("/view/:id", handler.ViewPost).Name("view")
	article.Post("/like/:id", handler.LikePost).Name("like")
}

func setupSharingRouter(r fiber.Router) {
	sharing := r.Group("/sharing").Name("sharing.")
	sharing.Post("/create", middleware.Auth, handler.CreateSharing)
	sharing.Delete("/delete/:id", middleware.Auth, handler.DeleteSharing)
	sharing.Put("/update", middleware.Auth, handler.UpdateSharing)
	sharing.Post("/list", handler.ListPost).Name("list")
	sharing.Get("/detail/:id", handler.DetailSharing).Name("detail")
	sharing.Post("/view/:id", handler.ViewPost).Name("view")
	sharing.Post("/like/:id", handler.LikePost).Name("like")
}

func setupMessageRouter(r fiber.Router) {
	message := r.Group("/message").Name("message.")
	message.Post("/create", middleware.Auth, handler.CreateMessage)
	message.Delete("/delete/:id", middleware.Auth, handler.DeleteMessage)
	message.Put("/update", middleware.Auth, handler.UpdateMessage)
	message.Post("/list", handler.ListMessage).Name("list")
	message.Get("/detail/:id", handler.DetailMessage).Name("detail")
	message.Post("/view/:id", handler.ViewMessage).Name("view")
	message.Post("/like/:id", handler.LikeMessage).Name("like")
}
