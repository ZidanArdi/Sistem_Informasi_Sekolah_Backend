package route

import (
	"backend/middleware"
	"backend/modules/school/handler"

	"github.com/gofiber/fiber/v2"
)

func SchoolRoute(router fiber.Router) {
	school := router.Group("/school-profile")
	school.Get("/", handler.GetSchoolProfile)
	school.Put("/", middleware.RequireAdmin, handler.UpdateSchoolProfile)
}
