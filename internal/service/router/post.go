package router

import (
	"Mou1ght/internal/api/handler"
	"Mou1ght/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func setupArticleRouter(r fiber.Router, articleHandler *handler.ArticleHandler, postHandler *handler.PostHandler) {
	article := r.Group("/article").Name("article.")
	article.Post("/create", middleware.Auth, articleHandler.CreateArticle)
	article.Delete("/delete/:id", middleware.Auth, articleHandler.DeleteArticle)
	article.Put("/update", middleware.Auth, articleHandler.UpdateArticle)
	article.Post("/status", middleware.Auth, postHandler.UpdatePostStatus)
	article.Post("/list", postHandler.ListPostPublic).Name("list")
	article.Post("/list/admin", middleware.Auth, postHandler.ListPost).Name("list_admin")
	article.Get("/detail/:id", articleHandler.DetailArticle).Name("detail")
	article.Post("/view/:id", postHandler.ViewPost).Name("view")
	article.Post("/like/:id", postHandler.LikePost).Name("like")
}

func setupSharingRouter(r fiber.Router, sharingHandler *handler.SharingHandler, postHandler *handler.PostHandler) {
	sharing := r.Group("/sharing").Name("sharing.")
	sharing.Post("/create", middleware.Auth, sharingHandler.CreateSharing)
	sharing.Delete("/delete/:id", middleware.Auth, sharingHandler.DeleteSharing)
	sharing.Put("/update", middleware.Auth, sharingHandler.UpdateSharing)
	sharing.Post("/list", postHandler.ListPostPublic).Name("list")
	sharing.Post("/list/admin", middleware.Auth, postHandler.ListPost).Name("list_admin")
	sharing.Post("/status", middleware.Auth, postHandler.UpdatePostStatus)
	sharing.Get("/detail/:id", sharingHandler.DetailSharing).Name("detail")
	sharing.Post("/view/:id", postHandler.ViewPost).Name("view")
	sharing.Post("/like/:id", postHandler.LikePost).Name("like")
}

func setupMessageRouter(r fiber.Router, messageHandler *handler.MessageHandler, postHandler *handler.PostHandler) {
	message := r.Group("/message").Name("message.")
	message.Get("/visitor", messageHandler.VisitorID)
	message.Post("/create", messageHandler.CreateMessage)
	message.Delete("/delete/:id", middleware.Auth, messageHandler.DeleteMessage)
	message.Put("/update", middleware.Auth, messageHandler.UpdateMessage)
	message.Post("/position", messageHandler.UpdateMessagePosition)
	message.Get("/owned", messageHandler.OwnedMessageIDs)
	message.Post("/list", messageHandler.ListMessagePublic).Name("list")
	message.Post("/list/admin", middleware.Auth, messageHandler.ListMessage).Name("list_admin")
	message.Post("/status", middleware.Auth, postHandler.UpdatePostStatus)
	message.Get("/detail/:id", messageHandler.DetailMessage).Name("detail")
	message.Post("/view/:id", messageHandler.ViewMessage).Name("view")
	message.Post("/like/:id", messageHandler.LikeMessage).Name("like")
}
