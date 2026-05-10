package router

import (
	"Mou1ght/internal/api/handler"
	"Mou1ght/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func setupTagRouter(r fiber.Router, tagHandler *handler.TagHandler) {
	tag := r.Group("/tag")
	tag.Use(middleware.Auth).Post("/create", tagHandler.CreateTag)
	tag.Use(middleware.Auth).Delete("/delete/:id", tagHandler.DeleteTag)
	tag.Get("/all", tagHandler.GetAllTags)
}

func setupCategoryRouter(r fiber.Router, categoryHandler *handler.CategoryHandler) {
	category := r.Group("/category")
	category.Use(middleware.Auth).Post("/create", categoryHandler.CreateCategory)
	category.Use(middleware.Auth).Put("/update/:id", categoryHandler.UpdateCategory)
	category.Use(middleware.Auth).Delete("/delete/:id", categoryHandler.DeleteCategory)
	category.Get("/all", categoryHandler.GetAllCategories)
}
