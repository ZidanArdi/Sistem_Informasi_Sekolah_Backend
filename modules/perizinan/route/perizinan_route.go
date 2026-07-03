package route

import (
	"backend/modules/perizinan/handler"

	"github.com/gofiber/fiber/v2"
)

func PerizinanRoute(app fiber.Router) {
	perizinan := app.Group("/perizinan")

	perizinan.Post("/", handler.CreatePerizinan)
	perizinan.Get("/", handler.GetPerizinan)
	perizinan.Put("/:id/approve", handler.ApproveOrRejectPerizinan)
}
