package router

import (
	"Mou1ght/consts"
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

func InitRouter() *fiber.App {
	r := fiber.New()
	r.Use(cors.New())
	setupRouter(r)
	return r
}

func setupRouter(r *fiber.App) {
	r.Static("/upload", consts.Upload, fiber.Static{
		// Download: true,
	})
	r.Use("/admin", filesystem.New(filesystem.Config{
		Root:         http.FS(adminUIFS()),
		NotFoundFile: "index.html",
	}))
	setupApiRouter(r)
}

func setupApiRouter(r *fiber.App) {
	v1 := r.Group("/api/v1")
	setupAttachmentRouter(v1)
	setupUserRouter(v1)
	setupSettingRouter(v1)
	setupArticleRouter(v1)
	setupSharingRouter(v1)
	setupMessageRouter(v1)
	setupTagRouter(v1)
	setupCategoryRouter(v1)
}
