package route

import (
	"backend/middleware"
	"backend/modules/guru/handler"

	"github.com/gofiber/fiber/v2"
)

func GuruRoute(app fiber.Router) {
	guru := app.Group("/guru")

	guru.Get("/", handler.GetAllGuru)
	guru.Get("/:id", handler.GetGuruByID)
	guru.Post("/", handler.CreateGuru)
	guru.Put("/:id", handler.UpdateGuru)
	guru.Delete("/:id", middleware.RequireAdmin, handler.DeleteGuru)
}
