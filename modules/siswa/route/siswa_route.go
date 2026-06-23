package route

import (
	"backend/middleware"
	"backend/modules/siswa/handler"

	"github.com/gofiber/fiber/v2"
)

func SiswaRoute(app fiber.Router) {

	siswa := app.Group("/siswa")

	siswa.Get("/", handler.GetAllSiswa)

	siswa.Get("/:id", handler.GetSiswaByID)

	siswa.Post("/", middleware.JWTProtected, handler.CreateSiswa)

	siswa.Put("/:id", middleware.JWTProtected, handler.UpdateSiswa)

	siswa.Delete("/:id", middleware.JWTProtected, middleware.RequireAdmin, handler.DeleteSiswa)
}
