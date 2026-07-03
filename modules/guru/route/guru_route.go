package route

import (
	"backend/middleware"
	authHandler "backend/modules/auth/handler"
	"backend/modules/guru/handler"

	"github.com/gofiber/fiber/v2"
)

func GuruRoute(app fiber.Router) {
	guru := app.Group("/guru")

	guru.Get("/", handler.GetAllGuru)
	guru.Get("/:id<int>", handler.GetGuruByID)
	guru.Post("/", handler.CreateGuru)
	guru.Put("/:id<int>", handler.UpdateGuru)
	guru.Delete("/:id<int>", middleware.RequireAdmin, handler.DeleteGuru)
	guru.Put("/change-password", authHandler.ChangePassword)
}
