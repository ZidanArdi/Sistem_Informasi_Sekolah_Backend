package handler

import (
	"backend/helpers"
	"backend/modules/auth/service"

	"github.com/gofiber/fiber/v2"
)

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
