package route

import (
	"backend/middleware"
	"backend/modules/nilai/handler"

	"github.com/gofiber/fiber/v2"
)

func NilaiRoute(app fiber.Router) {
	nilai := app.Group("/nilai")

	nilai.Get("/", handler.GetAllNilai)
	nilai.Get("/:id", handler.GetNilaiByID)
	nilai.Post("/", handler.CreateNilai)
	nilai.Put("/:id", handler.UpdateNilai)
	nilai.Delete("/:id", middleware.RequireAdmin, handler.DeleteNilai)
}
