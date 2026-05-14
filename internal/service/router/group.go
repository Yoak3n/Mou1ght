package router

import (
	"Mou1ght/internal/api/handler"
	"Mou1ght/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func setupTagRouter(r fiber.Router, tagHandler *handler.TagHandler) {
	tag := r.Group("/tag")
	tag.Post("/create", middleware.Auth, tagHandler.CreateTag)
	tag.Delete("/delete/:id", middleware.Auth, tagHandler.DeleteTag)
	tag.Get("/all", tagHandler.GetAllTags)
}

func setupCategoryRouter(r fiber.Router, categoryHandler *handler.CategoryHandler) {
	category := r.Group("/category")
	category.Post("/create", middleware.Auth, categoryHandler.CreateCategory)
	category.Put("/update/:id", middleware.Auth, categoryHandler.UpdateCategory)
	category.Delete("/delete/:id", middleware.Auth, categoryHandler.DeleteCategory)
	category.Get("/all", categoryHandler.GetAllCategories)
}
