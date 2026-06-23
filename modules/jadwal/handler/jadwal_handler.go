package handler

import (
	"strconv"

	"backend/helpers"
	"backend/modules/jadwal/model"
	"backend/modules/jadwal/repository"
	"backend/modules/jadwal/service"

	"github.com/gofiber/fiber/v2"
)

func GetAllJadwal(c *fiber.Ctx) error {
	data, err := service.GetAllJadwal(c.Query("kelas_id"), c.Query("mapel_id"), c.Query("guru_id"), c.Query("hari"))
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal mengambil data jadwal")
	}
	return helpers.SuccessResponse(c, "Berhasil mengambil data jadwal", data)
}

func GetJadwalByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID jadwal tidak valid")
	}

	data, err := service.GetJadwalByID(uint(id))
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Jadwal tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 500, "Gagal mengambil jadwal")
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil detail jadwal", data)
}

func CreateJadwal(c *fiber.Ctx) error {
	var jadwal model.Jadwal
	if err := c.BodyParser(&jadwal); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, err := service.CreateJadwal(jadwal)
	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Jadwal berhasil ditambahkan",
		"data":    data,
	})
}

func UpdateJadwal(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID jadwal tidak valid")
	}

	var jadwal model.Jadwal
	if err := c.BodyParser(&jadwal); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, err := service.UpdateJadwal(uint(id), jadwal)
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Jadwal tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return helpers.SuccessResponse(c, "Jadwal berhasil diupdate", data)
}

func DeleteJadwal(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID jadwal tidak valid")
	}

	err = service.DeleteJadwal(uint(id))
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Jadwal tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 500, "Gagal menghapus jadwal")
	}

	return helpers.SuccessResponse(c, "Jadwal berhasil dihapus", nil)
}
