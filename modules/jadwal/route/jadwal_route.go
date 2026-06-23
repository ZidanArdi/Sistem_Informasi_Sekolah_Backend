package route

import (
	"backend/middleware"
	"backend/modules/jadwal/handler"

	"github.com/gofiber/fiber/v2"
)

func JadwalRoute(app fiber.Router) {
	jadwal := app.Group("/jadwal")

	jadwal.Get("/", handler.GetAllJadwal)
	jadwal.Get("/:id", handler.GetJadwalByID)
	jadwal.Post("/", handler.CreateJadwal)
	jadwal.Put("/:id", handler.UpdateJadwal)
	jadwal.Delete("/:id", middleware.RequireAdmin, handler.DeleteJadwal)
}
