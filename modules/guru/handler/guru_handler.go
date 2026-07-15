package handler

import (
	"strconv"

	"backend/helpers"
	"backend/modules/guru/model"
	"backend/modules/guru/repository"
	"backend/modules/guru/service"

	"github.com/gofiber/fiber/v2"
)

// GetAllGuru godoc
// @Summary Ambil semua data guru
// @Description Mengambil seluruh data guru. Mendukung parameter pencarian.
// @Tags Guru
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Cari berdasarkan nama atau NIP"
// @Success 200 {object} model.GuruResponseList
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/guru/ [get]
func GetAllGuru(c *fiber.Ctx) error {
	data, err := service.GetAllGuru(c.Query("search"))
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal mengambil data guru")
	}
	return helpers.SuccessResponse(c, "Berhasil mengambil data guru", data)
}

// GetGuruByID godoc
// @Summary Ambil data guru berdasarkan ID
// @Description Mengambil detail data satu guru berdasarkan ID.
// @Tags Guru
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Guru"
// @Success 200 {object} model.GuruResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/guru/{id} [get]
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

// CreateGuru godoc
// @Summary Tambah data guru
// @Description Menambahkan data guru baru dan membuat user account dengan password acak sementara.
// @Tags Guru
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.Guru true "Payload data guru"
// @Success 201 {object} model.GuruResponseCreated
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/guru/ [post]
func CreateGuru(c *fiber.Ctx) error {
	var guru model.Guru
	if err := c.BodyParser(&guru); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, temporaryPassword, err := service.CreateGuru(guru)
	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"success":            true,
		"message":            "Guru berhasil ditambahkan",
		"data":               data,
		"temporary_password": temporaryPassword,
	})
}

// UpdateGuru godoc
// @Summary Ubah data guru
// @Description Mengubah data guru berdasarkan ID.
// @Tags Guru
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Guru"
// @Param request body model.Guru true "Payload data guru"
// @Success 200 {object} model.GuruResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/guru/{id} [put]
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

// DeleteGuru godoc
// @Summary Hapus data guru
// @Description Menghapus data guru beserta user account terkait berdasarkan ID. Hanya dapat diakses oleh Admin.
// @Tags Guru
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Guru"
// @Success 200 {object} helpers.SwaggerSuccessResponse
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 403 {object} helpers.SwaggerError403Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/guru/{id} [delete]
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
