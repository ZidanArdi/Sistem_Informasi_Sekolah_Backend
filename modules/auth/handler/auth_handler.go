package handler

import (
	"backend/helpers"
	authModel "backend/modules/auth/model"
	"backend/modules/auth/service"

	"github.com/gofiber/fiber/v2"
)

// Suppress unused import warning for Swagger docs references
var _ = authModel.User{}

// Register godoc
// @Summary Register user baru
// @Description Mendaftarkan akun user baru (khusus untuk role guru). Pendaftaran mandiri hanya diperbolehkan untuk role guru.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body service.RegisterInput true "Payload register user"
// @Success 201 {object} authModel.RegisterSuccessResponse
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/auth/register [post]
func Register(c *fiber.Ctx) error {
	var input service.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Input tidak valid")
	}

	user, err := service.Register(input)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Register berhasil",
		"data":    user,
	})
}

// Login godoc
// @Summary Login user
// @Description Melakukan autentikasi menggunakan Email/NIP/NIS/Username dan Password untuk mendapatkan JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body service.LoginInput true "Payload login user"
// @Success 200 {object} authModel.LoginSuccessResponse
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/auth/login [post]
func Login(c *fiber.Ctx) error {
	var input service.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Input tidak valid")
	}

	user, token, err := service.Login(input)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	return helpers.SuccessResponse(c, "Login berhasil", fiber.Map{
		"user":  user,
		"token": token,
	})
}

// ChangePassword godoc
// @Summary Ubah password user
// @Description Mengubah password user yang sedang login menggunakan token JWT.
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body service.ChangePasswordInput true "Payload ubah password"
// @Success 200 {object} helpers.SwaggerSuccessResponse
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 403 {object} helpers.SwaggerError403Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/auth/change-password [put]
func ChangePassword(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "User tidak valid")
	}

	var input service.ChangePasswordInput
	if err := c.BodyParser(&input); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Input tidak valid")
	}

	if err := service.ChangePassword(userID, input); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return helpers.SuccessResponse(c, "Password berhasil diubah", nil)
}
