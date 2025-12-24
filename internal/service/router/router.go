package router

import (
	"Mou1ght/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitRouter() *fiber.App {
	r := fiber.New()
	r.Use(cors.New())
	setupRouter(r)
	return r
}

func setupRouter(r *fiber.App) {
	r.Static("/upload", consts.Upload, fiber.Static{
		Download: true,
	})
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
