package handler

import (
	"strconv"

	"backend/helpers"
	"backend/modules/jadwal/model"
	"backend/modules/jadwal/repository"
	"backend/modules/jadwal/service"

	"github.com/gofiber/fiber/v2"
)

// GetAllJadwal godoc
// @Summary Ambil semua data jadwal
// @Description Mengambil seluruh data jadwal pelajaran dengan berbagai filter opsional.
// @Tags Jadwal
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param kelas_id query string false "Filter ID Kelas"
// @Param mapel_id query string false "Filter ID Mapel"
// @Param guru_id query string false "Filter ID Guru"
// @Param hari query string false "Filter Hari"
// @Param tahun_ajaran query string false "Filter Tahun Ajaran (e.g. 2026/2027)"
// @Param semester query string false "Filter Semester (e.g. Ganjil/Genap)"
// @Success 200 {object} model.JadwalResponseList
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/jadwal/ [get]
func GetAllJadwal(c *fiber.Ctx) error {
	data, err := service.GetAllJadwal(c.Query("kelas_id"), c.Query("mapel_id"), c.Query("guru_id"), c.Query("hari"), c.Query("tahun_ajaran"), c.Query("semester"))
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal mengambil data jadwal")
	}
	return helpers.SuccessResponse(c, "Berhasil mengambil data jadwal", data)
}

// GetJadwalByID godoc
// @Summary Ambil data jadwal berdasarkan ID
// @Description Mengambil detail data satu jadwal pelajaran berdasarkan ID.
// @Tags Jadwal
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Jadwal"
// @Success 200 {object} model.JadwalResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/jadwal/{id} [get]
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

// CreateJadwal godoc
// @Summary Tambah data jadwal
// @Description Menambahkan data jadwal pelajaran baru.
// @Tags Jadwal
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.Jadwal true "Payload data jadwal"
// @Success 201 {object} model.JadwalResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/jadwal/ [post]
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

// UpdateJadwal godoc
// @Summary Ubah data jadwal
// @Description Mengubah data jadwal pelajaran berdasarkan ID.
// @Tags Jadwal
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Jadwal"
// @Param request body model.Jadwal true "Payload data jadwal"
// @Success 200 {object} model.JadwalResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/jadwal/{id} [put]
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

// DeleteJadwal godoc
// @Summary Hapus data jadwal
// @Description Menghapus data jadwal pelajaran berdasarkan ID. Hanya dapat diakses oleh Admin.
// @Tags Jadwal
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Jadwal"
// @Success 200 {object} helpers.SwaggerSuccessResponse
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 403 {object} helpers.SwaggerError403Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/jadwal/{id} [delete]
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
