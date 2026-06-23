package route

import (
	"backend/middleware"
	"backend/modules/mapel/handler"

	"github.com/gofiber/fiber/v2"
)

func MapelRoute(app fiber.Router) {
	mapel := app.Group("/mapel")

	mapel.Get("/", handler.GetAllMapel)
	mapel.Get("/:id", handler.GetMapelByID)
	mapel.Post("/", handler.CreateMapel)
	mapel.Put("/:id", handler.UpdateMapel)
	mapel.Delete("/:id", middleware.RequireAdmin, handler.DeleteMapel)
}
