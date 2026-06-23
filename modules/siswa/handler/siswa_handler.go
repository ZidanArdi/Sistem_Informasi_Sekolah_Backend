package handler

import (
	"strconv"

	"backend/helpers"
	"backend/modules/siswa/model"
	"backend/modules/siswa/repository"
	"backend/modules/siswa/service"

	"github.com/gofiber/fiber/v2"
)

func GetAllSiswa(c *fiber.Ctx) error {

	search := c.Query("search")
	kelasID := c.Query("kelas_id")

	data, err := service.GetAllSiswa(search, kelasID)

	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal mengambil data siswa")
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil data siswa", data)
}

func GetSiswaByID(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID siswa tidak valid")
	}

	data, err := service.GetSiswaByID(uint(id))

	if err != nil {

		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Siswa tidak ditemukan")
		}

		return helpers.ErrorResponse(c, 500, "Gagal mengambil siswa")
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil detail siswa", data)
}

func CreateSiswa(c *fiber.Ctx) error {

	var siswa model.Siswa

	err := c.BodyParser(&siswa)

	if err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, err := service.CreateSiswa(siswa)

	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Siswa berhasil ditambahkan",
		"data":    data,
	})
}

func UpdateSiswa(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID siswa tidak valid")
	}

	var siswa model.Siswa

	err = c.BodyParser(&siswa)

	if err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, err := service.UpdateSiswa(uint(id), siswa)

	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Siswa tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return helpers.SuccessResponse(c, "Siswa berhasil diupdate", data)
}

func DeleteSiswa(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID siswa tidak valid")
	}

	err = service.DeleteSiswa(uint(id))

	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Siswa tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 500, "Gagal menghapus siswa")
	}

	return helpers.SuccessResponse(c, "Siswa berhasil dihapus", nil)
}
