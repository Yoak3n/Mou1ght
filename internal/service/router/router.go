package router

import (
	"Mou1ght/consts"
	"Mou1ght/internal/api/handler"
	"embed"
	"io/fs"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed all:adminui/*
var adminEmbedded embed.FS

func adminUIFS() fs.FS {
	if sub, err := fs.Sub(adminEmbedded, "adminui/dist"); err == nil {
		return sub
	}
	sub, err := fs.Sub(adminEmbedded, "adminui")
	if err != nil {
		panic(err)
	}
	return sub
}

type Deps struct {
	UserHandler       *handler.UserHandler
	ArticleHandler    *handler.ArticleHandler
	SharingHandler    *handler.SharingHandler
	MessageHandler    *handler.MessageHandler
	AttachmentHandler *handler.AttachmentHandler
	TagHandler        *handler.TagHandler
	CategoryHandler   *handler.CategoryHandler
	PostHandler       *handler.PostHandler
}

func InitRouter(deps Deps) *fiber.App {
	r := fiber.New()
	r.Use(cors.New())
	setupRouter(r, &deps)
	return r
}

func setupRouter(r *fiber.App, deps *Deps) {
	r.Static("/upload", consts.Upload, fiber.Static{
		// Download: true,
	})
	r.Use("/admin", filesystem.New(filesystem.Config{
		Root:         http.FS(adminUIFS()),
		NotFoundFile: "index.html",
	}))
	setupApiRouter(r, deps)
}

func setupApiRouter(r *fiber.App, deps *Deps) {
	v1 := r.Group("/api/v1")
	setupAttachmentRouter(v1, deps.AttachmentHandler)
	setupUserRouter(v1, deps.UserHandler)
	setupSettingRouter(v1)
	setupArticleRouter(v1, deps.ArticleHandler, deps.PostHandler)
	setupSharingRouter(v1, deps.SharingHandler, deps.PostHandler)
	setupMessageRouter(v1, deps.MessageHandler, deps.PostHandler)
	setupTagRouter(v1, deps.TagHandler)
	setupCategoryRouter(v1, deps.CategoryHandler)
}
