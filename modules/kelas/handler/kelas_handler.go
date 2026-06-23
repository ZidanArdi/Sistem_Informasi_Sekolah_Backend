package handler

import (
	"strconv"

	"backend/helpers"
	"backend/modules/kelas/model"
	"backend/modules/kelas/repository"
	"backend/modules/kelas/service"

	"github.com/gofiber/fiber/v2"
)

func GetAllKelas(c *fiber.Ctx) error {
	data, err := service.GetAllKelas(c.Query("search"), c.Query("tingkat"))
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal mengambil data kelas")
	}
	return helpers.SuccessResponse(c, "Berhasil mengambil data kelas", data)
}

func GetKelasByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID kelas tidak valid")
	}

	data, err := service.GetKelasByID(uint(id))
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Kelas tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 500, "Gagal mengambil kelas")
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil detail kelas", data)
}

func CreateKelas(c *fiber.Ctx) error {
	var kelas model.Kelas
	if err := c.BodyParser(&kelas); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, err := service.CreateKelas(kelas)
	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Kelas berhasil ditambahkan",
		"data":    data,
	})
}

func UpdateKelas(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID kelas tidak valid")
	}

	var kelas model.Kelas
	if err := c.BodyParser(&kelas); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, err := service.UpdateKelas(uint(id), kelas)
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Kelas tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return helpers.SuccessResponse(c, "Kelas berhasil diupdate", data)
}

func DeleteKelas(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID kelas tidak valid")
	}

	err = service.DeleteKelas(uint(id))
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Kelas tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 500, "Gagal menghapus kelas")
	}

	return helpers.SuccessResponse(c, "Kelas berhasil dihapus", nil)
}
