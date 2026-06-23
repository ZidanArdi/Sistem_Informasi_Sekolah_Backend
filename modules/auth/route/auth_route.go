package route

import (
	"backend/middleware"
	"backend/modules/auth/handler"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(app *fiber.App) {
	auth := app.Group("/api/auth")

	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)
	auth.Put("/change-password", middleware.JWTProtected, handler.ChangePassword)
}
