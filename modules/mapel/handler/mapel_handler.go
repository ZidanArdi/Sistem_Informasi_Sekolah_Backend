package handler

import (
	"strconv"

	"backend/helpers"
	"backend/modules/mapel/model"
	"backend/modules/mapel/repository"
	"backend/modules/mapel/service"

	"github.com/gofiber/fiber/v2"
)

func GetAllMapel(c *fiber.Ctx) error {
	data, err := service.GetAllMapel(c.Query("search"))
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal mengambil data mapel")
	}
	return helpers.SuccessResponse(c, "Berhasil mengambil data mapel", data)
}

func GetMapelByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID mapel tidak valid")
	}

	data, err := service.GetMapelByID(uint(id))
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Mapel tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 500, "Gagal mengambil mapel")
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil detail mapel", data)
}

func CreateMapel(c *fiber.Ctx) error {
	var mapel model.Mapel
	if err := c.BodyParser(&mapel); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, err := service.CreateMapel(mapel)
	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Mapel berhasil ditambahkan",
		"data":    data,
	})
}

func UpdateMapel(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID mapel tidak valid")
	}

	var mapel model.Mapel
	if err := c.BodyParser(&mapel); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, err := service.UpdateMapel(uint(id), mapel)
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Mapel tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return helpers.SuccessResponse(c, "Mapel berhasil diupdate", data)
}

func DeleteMapel(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID mapel tidak valid")
	}

	err = service.DeleteMapel(uint(id))
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Mapel tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 500, "Gagal menghapus mapel")
	}

	return helpers.SuccessResponse(c, "Mapel berhasil dihapus", nil)
}
