package route

import (
	"backend/middleware"
	"backend/modules/user/handler"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(router fiber.Router) {
	profile := router.Group("/profile")
	profile.Get("/", handler.GetProfile)
	profile.Put("/", handler.UpdateProfile)
	profile.Post("/photo", handler.UploadProfilePhoto)

	user := router.Group("/users", middleware.RequireAdmin)

	user.Get("/", handler.GetUsers)
	user.Put("/:id/status", handler.UpdateUserStatus)
	user.Put("/:id/reset-password", handler.ResetUserPassword)
}
