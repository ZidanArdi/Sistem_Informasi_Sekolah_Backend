package route

import (
	"backend/modules/siswa/handler"

	"github.com/gofiber/fiber/v2"
)

func SiswaRoute(app *fiber.App) {

	siswa := app.Group("/api/siswa")

	siswa.Get("/", handler.GetAllSiswa)

	siswa.Get("/:id", handler.GetSiswaByID)

	siswa.Post("/", handler.CreateSiswa)

	siswa.Put("/:id", handler.UpdateSiswa)

	siswa.Delete("/:id", handler.DeleteSiswa)
}