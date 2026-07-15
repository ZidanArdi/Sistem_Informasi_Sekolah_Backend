package handler

import (
	"strconv"

	"backend/helpers"
	"backend/modules/mapel/model"
	"backend/modules/mapel/repository"
	"backend/modules/mapel/service"

	"github.com/gofiber/fiber/v2"
)

// GetAllMapel godoc
// @Summary Ambil semua data mapel
// @Description Mengambil seluruh data mata pelajaran. Mendukung pencarian.
// @Tags Mapel
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Cari berdasarkan nama mapel atau kode"
// @Success 200 {object} model.MapelResponseList
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/mapel/ [get]
func GetAllMapel(c *fiber.Ctx) error {
	data, err := service.GetAllMapel(c.Query("search"))
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal mengambil data mapel")
	}
	return helpers.SuccessResponse(c, "Berhasil mengambil data mapel", data)
}

// GetMapelByID godoc
// @Summary Ambil data mapel berdasarkan ID
// @Description Mengambil detail data satu mata pelajaran berdasarkan ID.
// @Tags Mapel
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Mapel"
// @Success 200 {object} model.MapelResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/mapel/{id} [get]
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

// CreateMapel godoc
// @Summary Tambah data mapel
// @Description Menambahkan data mata pelajaran baru.
// @Tags Mapel
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.Mapel true "Payload data mapel"
// @Success 201 {object} model.MapelResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/mapel/ [post]
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

// UpdateMapel godoc
// @Summary Ubah data mapel
// @Description Mengubah data mata pelajaran berdasarkan ID.
// @Tags Mapel
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Mapel"
// @Param request body model.Mapel true "Payload data mapel"
// @Success 200 {object} model.MapelResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/mapel/{id} [put]
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

// DeleteMapel godoc
// @Summary Hapus data mapel
// @Description Menghapus data mata pelajaran berdasarkan ID. Hanya dapat diakses oleh Admin.
// @Tags Mapel
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Mapel"
// @Success 200 {object} helpers.SwaggerSuccessResponse
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 403 {object} helpers.SwaggerError403Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/mapel/{id} [delete]
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
