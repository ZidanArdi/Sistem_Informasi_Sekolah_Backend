package handler

import (
	"strconv"

	"backend/helpers"
	"backend/modules/guru/model"
	"backend/modules/guru/repository"
	"backend/modules/guru/service"

	"github.com/gofiber/fiber/v2"
)

func GetAllGuru(c *fiber.Ctx) error {
	data, err := service.GetAllGuru(c.Query("search"))
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal mengambil data guru")
	}
	return helpers.SuccessResponse(c, "Berhasil mengambil data guru", data)
}

func GetGuruByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID guru tidak valid")
	}

	data, err := service.GetGuruByID(uint(id))
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Guru tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 500, "Gagal mengambil guru")
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil detail guru", data)
}

func CreateGuru(c *fiber.Ctx) error {
	var guru model.Guru
	if err := c.BodyParser(&guru); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, err := service.CreateGuru(guru)
	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Guru berhasil ditambahkan",
		"data":    data,
	})
}

func UpdateGuru(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID guru tidak valid")
	}

	var guru model.Guru
	if err := c.BodyParser(&guru); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, err := service.UpdateGuru(uint(id), guru)
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Guru tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return helpers.SuccessResponse(c, "Guru berhasil diupdate", data)
}

func DeleteGuru(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID guru tidak valid")
	}

	err = service.DeleteGuru(uint(id))
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Guru tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 500, "Gagal menghapus guru")
	}

	return helpers.SuccessResponse(c, "Guru berhasil dihapus", nil)
}
