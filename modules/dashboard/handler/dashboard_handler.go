package handler

import (
	"strings"

	"backend/helpers"
	academicModel "backend/modules/academic/model"
	"backend/modules/dashboard/service"

	"github.com/gofiber/fiber/v2"
)

func sendError(c *fiber.Ctx, defaultCode int, err error) error {
	msg := err.Error()
	bizCode := ""
	if strings.Contains(msg, ": ") {
		parts := strings.SplitN(msg, ": ", 2)
		if strings.HasPrefix(parts[0], "ERR_") {
			bizCode = parts[0]
			msg = parts[1]
		}
	}

	code := defaultCode
	if bizCode == "ERR_FORBIDDEN" {
		code = fiber.StatusForbidden
	} else if bizCode == "ERR_PROFILE_NOT_FOUND" {
		code = fiber.StatusNotFound
	} else if bizCode == "ERR_INTERNAL_SERVER" {
		code = fiber.StatusInternalServerError
	}

	if bizCode == "" {
		bizCode = "ERR_INTERNAL_SERVER"
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"code":    bizCode,
		"message": msg,
		"data":    nil,
	})
}

// GetAdminDashboard godoc
// @Summary Ambil data dashboard admin
// @Description Mengambil statistik ringkas dashboard untuk role Admin (e.g. jumlah guru, siswa, kelas, mapel).
// @Tags Dashboard
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} helpers.SwaggerSuccessResponse
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 403 {object} helpers.SwaggerError403Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/admin/dashboard [get]
func GetAdminDashboard(c *fiber.Ctx) error {
	userIDLocal := c.Locals("user_id")
	if userIDLocal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    "ERR_UNAUTHORIZED",
			"message": "Token tidak ditemukan atau tidak valid",
			"data":    nil,
		})
	}
	userID, ok := userIDLocal.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    "ERR_UNAUTHORIZED",
			"message": "Format UserID dalam token tidak valid",
			"data":    nil,
		})
	}

	data, err := service.GetAdminDashboard(userID)
	if err != nil {
		return sendError(c, fiber.StatusForbidden, err)
	}

	return helpers.SuccessResponse(c, "Berhasil memuat dasbor admin", data)
}

// GetGuruDashboard godoc
// @Summary Ambil data dashboard guru
// @Description Mengambil statistik ringkas dashboard untuk Guru (e.g. jadwal mengajar, permohonan izin pending).
// @Tags Dashboard
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} helpers.SwaggerSuccessResponse
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 403 {object} helpers.SwaggerError403Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/guru/dashboard [get]
func GetGuruDashboard(c *fiber.Ctx) error {
	teacherCtxLocal := c.Locals("teacher_context")
	if teacherCtxLocal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    "ERR_UNAUTHORIZED",
			"message": "Token tidak valid atau akses ditolak",
			"data":    nil,
		})
	}

	teacherCtxMap, ok := teacherCtxLocal.(map[string]interface{})
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    "ERR_UNAUTHORIZED",
			"message": "Format Teacher Context tidak valid",
			"data":    nil,
		})
	}

	teacherCtx := academicModel.TeacherContext{
		UserID: teacherCtxMap["UserID"].(uint),
		GuruID: teacherCtxMap["GuruID"].(uint),
		Role:   teacherCtxMap["Role"].(string),
	}

	data, err := service.GetGuruDashboard(teacherCtx)
	if err != nil {
		return sendError(c, fiber.StatusForbidden, err)
	}

	return helpers.SuccessResponse(c, "Berhasil memuat dasbor guru", data)
}

// GetSiswaDashboard godoc
// @Summary Ambil data dashboard siswa
// @Description Mengambil statistik ringkas dashboard untuk Siswa (e.g. kelas, nilai rata-rata, absensi, jadwal hari ini).
// @Tags Dashboard
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} helpers.SwaggerSuccessResponse
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 403 {object} helpers.SwaggerError403Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/siswa/dashboard [get]
func GetSiswaDashboard(c *fiber.Ctx) error {
	userIDLocal := c.Locals("user_id")
	if userIDLocal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    "ERR_UNAUTHORIZED",
			"message": "Token tidak ditemukan atau tidak valid",
			"data":    nil,
		})
	}
	userID, ok := userIDLocal.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    "ERR_UNAUTHORIZED",
			"message": "Format UserID dalam token tidak valid",
			"data":    nil,
		})
	}

	data, err := service.GetSiswaDashboard(userID)
	if err != nil {
		return sendError(c, fiber.StatusForbidden, err)
	}

	return helpers.SuccessResponse(c, "Berhasil memuat dasbor siswa", data)
}
