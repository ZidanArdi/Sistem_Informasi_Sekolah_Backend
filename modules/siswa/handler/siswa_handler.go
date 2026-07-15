package handler

import (
	"strconv"

	"backend/helpers"
	academicModel "backend/modules/academic/model"
	"backend/modules/siswa/model"
	"backend/modules/siswa/repository"
	"backend/modules/siswa/service"

	"github.com/gofiber/fiber/v2"
)

// GetAllSiswa godoc
// @Summary Ambil semua data siswa
// @Description Mengambil seluruh data siswa. Mendukung parameter pencarian dan filter kelas. Route ini dapat diakses secara publik, namun Guru/Wali Kelas akan mendapatkan data terbatas sesuai kelas masing-masing jika login.
// @Tags Siswa
// @Accept json
// @Produce json
// @Param search query string false "Cari berdasarkan nama atau NIS"
// @Param kelas_id query string false "Filter berdasarkan ID kelas"
// @Success 200 {object} model.SiswaResponseList
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/siswa [get]
func GetAllSiswa(c *fiber.Ctx) error {

	search := c.Query("search")
	kelasID := c.Query("kelas_id")
	
	var teacherCtx academicModel.TeacherContext
	if tc, ok := c.Locals("teacher_context").(map[string]interface{}); ok {
		teacherCtx = academicModel.TeacherContext{
			UserID: tc["UserID"].(uint),
			GuruID: tc["GuruID"].(uint),
			Role:   tc["Role"].(string),
		}
	} else {
		// Fallback for non-guru users or if context not set
		if role, ok := c.Locals("role").(string); ok {
			teacherCtx.Role = role
		}
	}

	data, err := service.GetAllSiswa(teacherCtx, search, kelasID)

	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal mengambil data siswa")
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil data siswa", data)
}

// GetSiswaByID godoc
// @Summary Ambil data siswa berdasarkan ID
// @Description Mengambil detail data satu siswa berdasarkan ID.
// @Tags Siswa
// @Accept json
// @Produce json
// @Param id path int true "ID Siswa"
// @Success 200 {object} model.SiswaResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/siswa/{id} [get]
func GetSiswaByID(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID siswa tidak valid")
	}

	var teacherCtx academicModel.TeacherContext
	if tc, ok := c.Locals("teacher_context").(map[string]interface{}); ok {
		teacherCtx = academicModel.TeacherContext{
			UserID: tc["UserID"].(uint),
			GuruID: tc["GuruID"].(uint),
			Role:   tc["Role"].(string),
		}
	} else {
		if role, ok := c.Locals("role").(string); ok {
			teacherCtx.Role = role
		}
	}

	data, err := service.GetSiswaByID(teacherCtx, uint(id))

	if err != nil {

		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Siswa tidak ditemukan")
		}

		return helpers.ErrorResponse(c, 500, "Gagal mengambil siswa")
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil detail siswa", data)
}

// CreateSiswa godoc
// @Summary Tambah data siswa
// @Description Menambahkan data siswa baru dan membuat user account dengan password acak sementara. Hanya dapat diakses oleh user yang terautentikasi (JWT).
// @Tags Siswa
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.Siswa true "Payload data siswa"
// @Success 201 {object} model.SiswaResponseCreated
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/siswa [post]
func CreateSiswa(c *fiber.Ctx) error {

	var siswa model.Siswa

	err := c.BodyParser(&siswa)

	if err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, temporaryPassword, err := service.CreateSiswa(siswa)

	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"success":            true,
		"message":            "Siswa berhasil ditambahkan",
		"data":               data,
		"temporary_password": temporaryPassword,
	})
}

// UpdateSiswa godoc
// @Summary Ubah data siswa
// @Description Mengubah data siswa berdasarkan ID. Hanya dapat diakses oleh user yang terautentikasi (JWT).
// @Tags Siswa
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Siswa"
// @Param request body model.Siswa true "Payload data siswa"
// @Success 200 {object} model.SiswaResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/siswa/{id} [put]
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

// DeleteSiswa godoc
// @Summary Hapus data siswa
// @Description Menghapus data siswa beserta user account terkait berdasarkan ID. Hanya dapat diakses oleh Admin (JWT + Admin role).
// @Tags Siswa
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Siswa"
// @Success 200 {object} helpers.SwaggerSuccessResponse
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 403 {object} helpers.SwaggerError403Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/siswa/{id} [delete]
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
