package handler

import (
	"strconv"

	"backend/helpers"
	"backend/modules/nilai/model"
	"backend/modules/nilai/repository"
	"backend/modules/nilai/service"

	"github.com/gofiber/fiber/v2"
)

func GetAllNilai(c *fiber.Ctx) error {
	data, err := service.GetAllNilai(c.Query("siswa_id"), c.Query("mapel_id"), c.Query("semester"), c.Query("jenis_nilai"))
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal mengambil data nilai")
	}
	return helpers.SuccessResponse(c, "Berhasil mengambil data nilai", data)
}

func GetNilaiByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID nilai tidak valid")
	}

	data, err := service.GetNilaiByID(uint(id))
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Nilai tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 500, "Gagal mengambil nilai")
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil detail nilai", data)
}

func CreateNilai(c *fiber.Ctx) error {
	var nilai model.Nilai
	if err := c.BodyParser(&nilai); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, err := service.CreateNilai(nilai)
	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Nilai berhasil ditambahkan",
		"data":    data,
	})
}

func UpdateNilai(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID nilai tidak valid")
	}

	var nilai model.Nilai
	if err := c.BodyParser(&nilai); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, err := service.UpdateNilai(uint(id), nilai)
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Nilai tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return helpers.SuccessResponse(c, "Nilai berhasil diupdate", data)
}

func DeleteNilai(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID nilai tidak valid")
	}

	err = service.DeleteNilai(uint(id))
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Nilai tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 500, "Gagal menghapus nilai")
	}

	return helpers.SuccessResponse(c, "Nilai berhasil dihapus", nil)
}
