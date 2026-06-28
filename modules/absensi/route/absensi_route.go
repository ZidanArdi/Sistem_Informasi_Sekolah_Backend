package route

import (
	"backend/middleware"
	"backend/modules/absensi/handler"

	"github.com/gofiber/fiber/v2"
)

func AbsensiRoute(app fiber.Router) {
	absensi := app.Group("/absensi")

	// Protected routes under JWT
	absensi.Use(middleware.JWTProtected)

	absensi.Get("/", handler.GetAllAbsensi)
	absensi.Get("/:id", handler.GetAbsensiByID)
	absensi.Post("/", handler.CreateAbsensi)
	absensi.Post("/bulk", handler.BulkSaveAbsensi)
	absensi.Put("/:id/approve", handler.ApproveOrRejectAbsensi)
	absensi.Delete("/:id", middleware.RequireAdmin, handler.DeleteAbsensi)
}
