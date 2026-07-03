package route

import (
	"backend/modules/dashboard/handler"

	"github.com/gofiber/fiber/v2"
)

func DashboardRoute(app fiber.Router) {
	app.Get("/admin/dashboard", handler.GetAdminDashboard)
	app.Get("/guru/dashboard", handler.GetGuruDashboard)
	app.Get("/siswa/dashboard", handler.GetSiswaDashboard)
}
