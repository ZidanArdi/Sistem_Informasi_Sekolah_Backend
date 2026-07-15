package handler

import (
	"strconv"

	"backend/helpers"
	academicModel "backend/modules/academic/model"
	"backend/modules/nilai/model"
	"backend/modules/nilai/repository"
	"backend/modules/nilai/service"

	"github.com/gofiber/fiber/v2"
)

func getTeacherContext(c *fiber.Ctx) academicModel.TeacherContext {
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
		if userID, ok := c.Locals("user_id").(uint); ok {
			teacherCtx.UserID = userID
		}
	}
	return teacherCtx
}

// GetAllNilai godoc
// @Summary Ambil semua data nilai
// @Description Mengambil seluruh data nilai siswa berdasarkan filter kelas, mapel, siswa, dll. Guru hanya dapat melihat kelas/mapel yang diajarkannya, sedangkan siswa hanya dapat melihat nilainya sendiri.
// @Tags Nilai
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param kelas_id query string false "Filter ID Kelas"
// @Param mapel_id query string false "Filter ID Mapel"
// @Param siswa_id query string false "Filter ID Siswa"
// @Param tahun_ajaran query string false "Filter Tahun Ajaran (e.g. 2026/2027)"
// @Param semester query string false "Filter Semester (e.g. Ganjil/Genap)"
// @Success 200 {object} model.NilaiResponseList
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/nilai/ [get]
func GetAllNilai(c *fiber.Ctx) error {
	teacherCtx := getTeacherContext(c)
	if teacherCtx.UserID == 0 {
		return helpers.ErrorResponse(c, 401, "User tidak valid")
	}

	data, err := service.GetAllNilai(
		teacherCtx,
		c.Query("siswa_id"),
		c.Query("kelas_id"),
		c.Query("mapel_id"),
		c.Query("semester"),
		c.Query("tahun_ajaran"),
	)
	if err != nil {
		return helpers.ErrorResponse(c, 500, err.Error())
	}
	return helpers.SuccessResponse(c, "Berhasil mengambil data nilai", data)
}

// GetNilaiByID godoc
// @Summary Ambil data nilai berdasarkan ID
// @Description Mengambil detail data satu nilai berdasarkan ID.
// @Tags Nilai
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Nilai"
// @Success 200 {object} model.NilaiResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/nilai/{id} [get]
func GetNilaiByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID nilai tidak valid")
	}

	data, err := service.GetNilaiByID(uint(id))
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Nilai tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 500, "Gagal mengambil nilai")
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil detail nilai", data)
}

// CreateNilai godoc
// @Summary Tambah data nilai
// @Description Menambahkan data nilai baru untuk siswa. Hanya guru pengampu/wali kelas/admin yang diijinkan mengisi nilai.
// @Tags Nilai
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.Nilai true "Payload data nilai"
// @Success 201 {object} model.NilaiResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/nilai/ [post]
func CreateNilai(c *fiber.Ctx) error {
	teacherCtx := getTeacherContext(c)
	if teacherCtx.UserID == 0 {
		return helpers.ErrorResponse(c, 401, "User tidak valid")
	}

	var nilai model.Nilai
	if err := c.BodyParser(&nilai); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, err := service.CreateNilai(teacherCtx, nilai)
	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Nilai berhasil ditambahkan",
		"data":    data,
	})
}

// UpdateNilai godoc
// @Summary Ubah data nilai
// @Description Mengubah data nilai siswa berdasarkan ID. Hanya dapat diubah oleh Guru pengampu atau Admin.
// @Tags Nilai
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Nilai"
// @Param request body model.Nilai true "Payload data nilai"
// @Success 200 {object} model.NilaiResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/nilai/{id} [put]
func UpdateNilai(c *fiber.Ctx) error {
	teacherCtx := getTeacherContext(c)
	if teacherCtx.UserID == 0 {
		return helpers.ErrorResponse(c, 401, "User tidak valid")
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID nilai tidak valid")
	}

	var nilai model.Nilai
	if err := c.BodyParser(&nilai); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, err := service.UpdateNilai(uint(id), teacherCtx, nilai)
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Nilai tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return helpers.SuccessResponse(c, "Nilai berhasil diupdate", data)
}

// DeleteNilai godoc
// @Summary Hapus data nilai
// @Description Menghapus data nilai siswa berdasarkan ID. Hanya dapat dihapus oleh Admin (RequireAdmin).
// @Tags Nilai
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Nilai"
// @Success 200 {object} helpers.SwaggerSuccessResponse
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 403 {object} helpers.SwaggerError403Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/nilai/{id} [delete]
func DeleteNilai(c *fiber.Ctx) error {
	teacherCtx := getTeacherContext(c)
	if teacherCtx.UserID == 0 {
		return helpers.ErrorResponse(c, 401, "User tidak valid")
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID nilai tidak valid")
	}

	err = service.DeleteNilai(uint(id), teacherCtx)
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Nilai tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 500, err.Error())
	}

	return helpers.SuccessResponse(c, "Nilai berhasil dihapus", nil)
}
