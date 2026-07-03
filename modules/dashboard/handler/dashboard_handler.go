package handler

import (
	"backend/helpers"
	"backend/modules/dashboard/service"

	"github.com/gofiber/fiber/v2"
)

func GetAdminDashboard(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)

	data, err := service.GetAdminDashboard(role)
	if err != nil {
		return helpers.ErrorResponse(c, 403, err.Error())
	}

	return helpers.SuccessResponse(c, "Berhasil memuat dasbor admin", data)
}

func GetGuruDashboard(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)
	email, _ := c.Locals("email").(string)

	data, err := service.GetGuruDashboard(role, email)
	if err != nil {
		return helpers.ErrorResponse(c, 403, err.Error())
	}

	return helpers.SuccessResponse(c, "Berhasil memuat dasbor guru", data)
}

func GetSiswaDashboard(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)
	email, _ := c.Locals("email").(string)

	data, err := service.GetSiswaDashboard(role, email)
	if err != nil {
		return helpers.ErrorResponse(c, 403, err.Error())
	}

	return helpers.SuccessResponse(c, "Berhasil memuat dasbor siswa", data)
}
