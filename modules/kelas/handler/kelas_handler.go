package handler

import (
	"strconv"

	"backend/helpers"
	"backend/modules/kelas/model"
	"backend/modules/kelas/repository"
	"backend/modules/kelas/service"

	"github.com/gofiber/fiber/v2"
)

// GetAllKelas godoc
// @Summary Ambil semua data kelas
// @Description Mengambil seluruh data kelas. Mendukung pencarian nama kelas dan filter tingkat.
// @Tags Kelas
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Cari berdasarkan nama kelas"
// @Param tingkat query string false "Filter tingkat kelas (e.g. X, XI, XII)"
// @Success 200 {object} model.KelasResponseList
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/kelas/ [get]
func GetAllKelas(c *fiber.Ctx) error {
	data, err := service.GetAllKelas(c.Query("search"), c.Query("tingkat"))
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal mengambil data kelas")
	}
	return helpers.SuccessResponse(c, "Berhasil mengambil data kelas", data)
}

// GetKelasByID godoc
// @Summary Ambil data kelas berdasarkan ID
// @Description Mengambil detail data satu kelas berdasarkan ID.
// @Tags Kelas
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Kelas"
// @Success 200 {object} model.KelasResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/kelas/{id} [get]
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

// CreateKelas godoc
// @Summary Tambah data kelas
// @Description Menambahkan data kelas baru.
// @Tags Kelas
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.Kelas true "Payload data kelas"
// @Success 201 {object} model.KelasResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/kelas/ [post]
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

// UpdateKelas godoc
// @Summary Ubah data kelas
// @Description Mengubah data kelas berdasarkan ID.
// @Tags Kelas
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Kelas"
// @Param request body model.Kelas true "Payload data kelas"
// @Success 200 {object} model.KelasResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/kelas/{id} [put]
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

// DeleteKelas godoc
// @Summary Hapus data kelas
// @Description Menghapus data kelas berdasarkan ID. Hanya dapat diakses oleh Admin.
// @Tags Kelas
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Kelas"
// @Success 200 {object} helpers.SwaggerSuccessResponse
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 403 {object} helpers.SwaggerError403Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/kelas/{id} [delete]
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
