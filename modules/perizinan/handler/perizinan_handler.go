package handler

import (
	"strconv"

	"backend/helpers"
	"backend/modules/perizinan/model"
	"backend/modules/perizinan/repository"
	"backend/modules/perizinan/service"

	"github.com/gofiber/fiber/v2"
)

type ApproveRejectPerizinanInput struct {
	Status         string `json:"status" example:"Disetujui"`
	KeteranganGuru string `json:"keterangan_guru" example:"Izin disetujui, semoga lekas sembuh."`
}

// CreatePerizinan godoc
// @Summary Ajukan permohonan perizinan
// @Description Mengajukan perizinan baru (misal izin sakit/keperluan keluarga) untuk siswa.
// @Tags Perizinan
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.Perizinan true "Payload perizinan"
// @Success 201 {object} model.PerizinanResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/perizinan/ [post]
func CreatePerizinan(c *fiber.Ctx) error {
	var perizinan model.Perizinan
	if err := c.BodyParser(&perizinan); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	role, _ := c.Locals("role").(string)
	email, _ := c.Locals("email").(string)

	data, err := service.CreatePerizinan(perizinan, role, email)
	if err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Permohonan perizinan berhasil diajukan",
		"data":    data,
	})
}

// GetPerizinan godoc
// @Summary Ambil data perizinan
// @Description Mengambil seluruh data perizinan yang relevan bagi user yang sedang login (siswa hanya melihat miliknya, guru melihat data kelasnya/seluruhnya).
// @Tags Perizinan
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.PerizinanResponseList
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/perizinan/ [get]
func GetPerizinan(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)
	email, _ := c.Locals("email").(string)

	data, err := service.GetPerizinan(role, email)
	if err != nil {
		return helpers.ErrorResponse(c, 500, err.Error())
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil data perizinan", data)
}

// ApproveOrRejectPerizinan godoc
// @Summary Setujui atau tolak permohonan perizinan
// @Description Mengubah status persetujuan perizinan siswa. Hanya dapat dilakukan oleh Guru.
// @Tags Perizinan
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Perizinan"
// @Param request body ApproveRejectPerizinanInput true "Payload persetujuan"
// @Success 200 {object} model.PerizinanResponseSingle
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/perizinan/{id}/approve [put]
func ApproveOrRejectPerizinan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.ErrorResponse(c, 400, "ID tidak valid")
	}

	var input ApproveRejectPerizinanInput

	if err := c.BodyParser(&input); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	role, _ := c.Locals("role").(string)
	email, _ := c.Locals("email").(string)

	data, err := service.ApproveOrRejectPerizinan(uint(id), input.Status, input.KeteranganGuru, role, email)
	if err != nil {
		if repository.IsNotFoundError(err) {
			return helpers.ErrorResponse(c, 404, "Data perizinan tidak ditemukan")
		}
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	return helpers.SuccessResponse(c, "Status persetujuan perizinan berhasil diperbarui", data)
}
