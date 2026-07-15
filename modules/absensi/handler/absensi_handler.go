package handler

import (
	"strconv"

	"backend/helpers"
	"backend/modules/absensi/model"
	"backend/modules/absensi/repository"
	"backend/modules/absensi/service"

	"github.com/gofiber/fiber/v2"
)

type BulkSaveInput struct {
	Tanggal string          `json:"tanggal" example:"2026-07-15"`
	Records []model.Absensi `json:"records"`
}

type ApproveRejectInput struct {
	StatusPersetujuan string `json:"status_persetujuan" example:"Disetujui"`
}

// GetAllAbsensi godoc
// @Summary Ambil semua data absensi
// @Description Mengambil seluruh data absensi. Dapat disaring berdasarkan kelas, tanggal, status persetujuan, dan siswa. Route ini terproteksi jika diakses via route ber-auth, namun absensiRoute juga didaftarkan sebagai publicAPI di main.go.
// @Tags Absensi
// @Accept json
// @Produce json
// @Param kelas_id query string false "Filter ID Kelas"
// @Param tanggal query string false "Filter Tanggal (YYYY-MM-DD)"
// @Param status_persetujuan query string false "Filter Status Persetujuan (e.g. Pending/Disetujui/Ditolak)"
// @Param siswa_id query int false "Filter ID Siswa"
// @Success 200 {object} model.AbsensiResponseList
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/absensi/ [get]
func GetAllAbsensi(c *fiber.Ctx) error {
	kelasID := c.Query("kelas_id")
	tanggal := c.Query("tanggal")
	statusPersetujuan := c.Query("status_persetujuan")
	
	var specificSiswaID uint = 0
	if siswaIDStr := c.Query("siswa_id"); siswaIDStr != "" {
		if id, err := strconv.Atoi(siswaIDStr); err == nil {
			specificSiswaID = uint(id)
		}
	}

	role, _ := c.Locals("role").(string)
	email, _ := c.Locals("email").(string)

	data, err := service.GetAllAbsensi(kelasID, tanggal, statusPersetujuan, specificSiswaID, role, email)
	if err != nil {
		return helpers.ErrorResponse(c, 500, err.Error())
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil data absensi", data)
}

// GetAbsensiByID godoc
// @Summary Ambil data absensi berdasarkan ID
// @Description Mengambil detail data satu absensi berdasarkan ID.
// @Tags Absensi
// @Accept json
// @Produce json
// @Param id path int true "ID Absensi"
// @Success 200 {object} model.AbsensiResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/absensi/{id} [get]
func GetAbsensiByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID tidak valid")
	}

	data, err := service.GetAbsensiByID(uint(id))
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Data absensi tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 500, "Gagal mengambil data absensi")
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil detail absensi", data)
}

// CreateAbsensi godoc
// @Summary Tambah data absensi / ajukan izin
// @Description Menambahkan data absensi baru atau mengajukan izin.
// @Tags Absensi
// @Accept json
// @Produce json
// @Param request body model.Absensi true "Payload data absensi"
// @Success 201 {object} model.AbsensiResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/absensi/ [post]
func CreateAbsensi(c *fiber.Ctx) error {
	var absensi model.Absensi
	if err := c.BodyParser(&absensi); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	role, _ := c.Locals("role").(string)
	email, _ := c.Locals("email").(string)

	data, err := service.CreateAbsensi(absensi, role, email)
	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Absensi/permohonan izin berhasil diajukan",
		"data":    data,
	})
}

// BulkSaveAbsensi godoc
// @Summary Simpan absensi kelas secara bulk
// @Description Menyimpan atau memperbarui data absensi banyak siswa sekaligus di suatu tanggal.
// @Tags Absensi
// @Accept json
// @Produce json
// @Param request body BulkSaveInput true "Payload bulk absensi"
// @Success 200 {object} helpers.SwaggerSuccessResponse
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/absensi/bulk [post]
func BulkSaveAbsensi(c *fiber.Ctx) error {
	var input BulkSaveInput

	if err := c.BodyParser(&input); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	err := service.BulkSaveAbsensi(input.Records, input.Tanggal)
	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return helpers.SuccessResponse(c, "Absensi kelas berhasil disimpan", nil)
}

// ApproveOrRejectAbsensi godoc
// @Summary Setujui atau tolak izin absensi
// @Description Mengubah status persetujuan absensi (e.g. Disetujui/Ditolak).
// @Tags Absensi
// @Accept json
// @Produce json
// @Param id path int true "ID Absensi"
// @Param request body ApproveRejectInput true "Payload status persetujuan"
// @Success 200 {object} model.AbsensiResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/absensi/{id}/approve [put]
func ApproveOrRejectAbsensi(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID tidak valid")
	}

	var input ApproveRejectInput

	if err := c.BodyParser(&input); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	data, err := service.ApproveOrRejectAbsensi(uint(id), input.StatusPersetujuan)
	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return helpers.SuccessResponse(c, "Status persetujuan absensi berhasil diperbarui", data)
}

// DeleteAbsensi godoc
// @Summary Hapus data absensi
// @Description Menghapus data absensi berdasarkan ID.
// @Tags Absensi
// @Accept json
// @Produce json
// @Param id path int true "ID Absensi"
// @Success 200 {object} helpers.SwaggerSuccessResponse
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/absensi/{id} [delete]
func DeleteAbsensi(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID tidak valid")
	}

	err = service.DeleteAbsensi(uint(id))
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal menghapus data absensi")
	}

	return helpers.SuccessResponse(c, "Data absensi berhasil dihapus", nil)
}
