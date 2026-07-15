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
