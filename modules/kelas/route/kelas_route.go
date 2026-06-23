package route

import (
	"backend/middleware"
	"backend/modules/kelas/handler"

	"github.com/gofiber/fiber/v2"
)

func KelasRoute(app fiber.Router) {
	kelas := app.Group("/kelas")

	kelas.Get("/", handler.GetAllKelas)
	kelas.Get("/:id", handler.GetKelasByID)
	kelas.Post("/", handler.CreateKelas)
	kelas.Put("/:id", handler.UpdateKelas)
	kelas.Delete("/:id", middleware.RequireAdmin, handler.DeleteKelas)
}
