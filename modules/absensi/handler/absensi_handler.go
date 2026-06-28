package handler

import (
	"strconv"

	"backend/helpers"
	"backend/modules/absensi/model"
	"backend/modules/absensi/repository"
	"backend/modules/absensi/service"

	"github.com/gofiber/fiber/v2"
)

func GetAllAbsensi(c *fiber.Ctx) error {
	kelasID := c.Query("kelas_id")
	tanggal := c.Query("tanggal")
	statusPersetujuan := c.Query("status_persetujuan")
	
	var specificSiswaID uint = 0
	if siswaIDStr := c.Query("siswa_id"); siswaIDStr != "" {
		if id, err := strconv.Atoi(siswaIDStr); err == nil {
			specificSiswaID = uint(id)
		}
	}

	role, _ := c.Locals("role").(string)
	email, _ := c.Locals("email").(string)

	data, err := service.GetAllAbsensi(kelasID, tanggal, statusPersetujuan, specificSiswaID, role, email)
	if err != nil {
		return helpers.ErrorResponse(c, 500, err.Error())
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil data absensi", data)
}

func GetAbsensiByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID tidak valid")
	}

	data, err := service.GetAbsensiByID(uint(id))
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Data absensi tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 500, "Gagal mengambil data absensi")
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil detail absensi", data)
}

func CreateAbsensi(c *fiber.Ctx) error {
	var absensi model.Absensi
	if err := c.BodyParser(&absensi); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	role, _ := c.Locals("role").(string)
	email, _ := c.Locals("email").(string)

	data, err := service.CreateAbsensi(absensi, role, email)
	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Absensi/permohonan izin berhasil diajukan",
		"data":    data,
	})
}

func BulkSaveAbsensi(c *fiber.Ctx) error {
	var input struct {
		Tanggal string          `json:"tanggal"`
		Records []model.Absensi `json:"records"`
	}

	if err := c.BodyParser(&input); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	err := service.BulkSaveAbsensi(input.Records, input.Tanggal)
	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return helpers.SuccessResponse(c, "Absensi kelas berhasil disimpan", nil)
}

func ApproveOrRejectAbsensi(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID tidak valid")
	}

	var input struct {
		StatusPersetujuan string `json:"status_persetujuan"`
	}

	if err := c.BodyParser(&input); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, err := service.ApproveOrRejectAbsensi(uint(id), input.StatusPersetujuan)
	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return helpers.SuccessResponse(c, "Status persetujuan absensi berhasil diperbarui", data)
}

func DeleteAbsensi(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID tidak valid")
	}

	err = service.DeleteAbsensi(uint(id))
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal menghapus data absensi")
	}

	return helpers.SuccessResponse(c, "Data absensi berhasil dihapus", nil)
}
