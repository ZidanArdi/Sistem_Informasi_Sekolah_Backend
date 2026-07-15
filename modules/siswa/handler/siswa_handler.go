package handler

import (
	"strconv"

	"backend/helpers"
	"backend/modules/siswa/model"
	"backend/modules/siswa/repository"
	"backend/modules/siswa/service"

	"github.com/gofiber/fiber/v2"
	academicModel "backend/modules/academic/model"
)

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
