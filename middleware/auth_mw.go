package middleware

import (
	"strings"

	"backend/helpers"

	"github.com/gofiber/fiber/v2"
)

func JWTProtected(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Token tidak ditemukan")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Format token tidak valid")
	}

	claims, err := helpers.ParseToken(parts[1])
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	c.Locals("user_id", claims.UserID)
	c.Locals("email", claims.Email)
	c.Locals("role", claims.Role)

	return c.Next()
}

func RequireAdmin(c *fiber.Ctx) error {
	role, ok := c.Locals("role").(string)
	if !ok || role != "admin" {
		return helpers.ErrorResponse(c, fiber.StatusForbidden, "Akses hanya untuk admin")
	}

	return c.Next()
}
