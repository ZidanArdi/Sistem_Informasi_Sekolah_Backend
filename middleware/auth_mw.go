package middleware

import (
	"strings"

	"backend/config"
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

	// Enforce first login password change policy
	var isFirstLogin bool
	config.DB.Table("users").Select("is_first_login").Where("id = ?", claims.UserID).Row().Scan(&isFirstLogin)

	if isFirstLogin && c.Path() != "/api/auth/change-password" {
		return helpers.ErrorResponse(c, fiber.StatusForbidden, "ERR_FIRST_LOGIN: FORCE_PASSWORD_CHANGE: Anda wajib mengubah password default terlebih dahulu")
	}

	// Initialize Lightweight TeacherContext
	if claims.Role == "guru" {
		var guruID uint
		config.DB.Table("gurus").Select("id").Where("user_id = ? AND deleted_at IS NULL", claims.UserID).Row().Scan(&guruID)
		
		teacherCtx := map[string]interface{}{
			"UserID": claims.UserID,
			"GuruID": guruID,
			"Role":   claims.Role,
		}
		c.Locals("teacher_context", teacherCtx)
	}

	return c.Next()
}

func RequireAdmin(c *fiber.Ctx) error {
	role, ok := c.Locals("role").(string)
	if !ok || role != "admin" {
		return helpers.ErrorResponse(c, fiber.StatusForbidden, "ERR_PERMISSION_DENIED: Akses hanya untuk admin")
	}

	return c.Next()
}
