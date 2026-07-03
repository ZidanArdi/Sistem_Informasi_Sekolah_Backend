package handler

import (
	"strconv"

	"backend/helpers"
	"backend/modules/perizinan/model"
	"backend/modules/perizinan/repository"
	"backend/modules/perizinan/service"

	"github.com/gofiber/fiber/v2"
)

func CreatePerizinan(c *fiber.Ctx) error {
	var perizinan model.Perizinan
	if err := c.BodyParser(&perizinan); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	role, _ := c.Locals("role").(string)
	email, _ := c.Locals("email").(string)

	data, err := service.CreatePerizinan(perizinan, role, email)
	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Permohonan perizinan berhasil diajukan",
		"data":    data,
	})
}

func GetPerizinan(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)
	email, _ := c.Locals("email").(string)

	data, err := service.GetPerizinan(role, email)
	if err != nil {
		return helpers.ErrorResponse(c, 500, err.Error())
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil data perizinan", data)
}

func ApproveOrRejectPerizinan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID tidak valid")
	}

	var input struct {
		Status         string `json:"status"`
		KeteranganGuru string `json:"keterangan_guru"`
	}

	if err := c.BodyParser(&input); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	role, _ := c.Locals("role").(string)
	email, _ := c.Locals("email").(string)

	data, err := service.ApproveOrRejectPerizinan(uint(id), input.Status, input.KeteranganGuru, role, email)
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Data perizinan tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return helpers.SuccessResponse(c, "Status persetujuan perizinan berhasil diperbarui", data)
}
