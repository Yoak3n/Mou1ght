package router

import (
	"Mou1ght/internal/api/handler"
	"Mou1ght/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func setupArticleRouter(r fiber.Router) {
	article := r.Group("/article").Name("article.")
	article.Use(middleware.Auth).Post("/create", handler.CreateArticle)
	article.Use(middleware.Auth).Delete("/delete/:id", handler.DeleteArticle)
	article.Use(middleware.Auth).Put("/update", handler.UpdateArticle)
	article.Post("/list", handler.ListPost).Name("list")
	article.Get("/detail/:id", handler.DetailArticle).Name("detail")
	article.Post("/view/:id", handler.ViewPost).Name("view")
	article.Post("/like/:id", handler.LikePost).Name("like")
}

func setupSharingRouter(r fiber.Router) {
	sharing := r.Group("/sharing").Name("sharing.")
	sharing.Use(middleware.Auth).Post("/create", handler.CreateSharing)
	sharing.Use(middleware.Auth).Delete("/delete/:id", handler.DeleteSharing)
	sharing.Use(middleware.Auth).Put("/update", handler.UpdateSharing)
	sharing.Post("/list", handler.ListPost).Name("list")
	sharing.Get("/detail/:id", handler.DetailSharing).Name("detail")
	sharing.Post("/view/:id", handler.ViewPost).Name("view")
	sharing.Post("/like/:id", handler.LikePost).Name("like")
}

func setupMessageRouter(r fiber.Router) {
	message := r.Group("/message").Name("message.")
	message.Use(middleware.Auth).Post("/create", handler.CreateMessage)
	message.Use(middleware.Auth).Delete("/delete/:id", handler.DeleteMessage)
	message.Use(middleware.Auth).Put("/update", handler.UpdateMessage)
	message.Post("/list", handler.ListMessage).Name("list")
	message.Get("/detail/:id", handler.DetailMessage).Name("detail")
	message.Post("/view/:id", handler.ViewMessage).Name("view")
	message.Post("/like/:id", handler.LikeMessage).Name("like")
}
