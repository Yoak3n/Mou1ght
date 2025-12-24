package router

import (
	"Mou1ght/internal/api/handler"
	"Mou1ght/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func setupTagRouter(r fiber.Router) {
	tag := r.Group("/tag")
	tag.Use(middleware.Auth).Post("/create", handler.CreateTag)
	tag.Use(middleware.Auth).Delete("/delete/:id", handler.DeleteTag)
	tag.Get("/all", handler.GetAllTags)
}

func setupCategoryRouter(r fiber.Router) {
	category := r.Group("/category")
	category.Use(middleware.Auth).Post("/create", handler.CreateCategory)
	category.Use(middleware.Auth).Put("/update/:id", handler.UpdateCategory)
	category.Use(middleware.Auth).Delete("/delete/:id", handler.DeleteCategory)
	category.Get("/all", handler.GetAllCategories)
}
