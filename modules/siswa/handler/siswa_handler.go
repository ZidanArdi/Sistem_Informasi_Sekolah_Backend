package handler

import (
	"regexp"

	"backend/helpers"
	"backend/modules/siswa/model"
	"backend/modules/siswa/repository"
	"backend/modules/siswa/service"

	"github.com/gofiber/fiber/v2"
)

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func GetAllSiswa(c *fiber.Ctx) error {

	search := c.Query("search")

	data, err := service.GetAllSiswa(search)

	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal mengambil data siswa")
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil data siswa", data)
}

func GetSiswaByID(c *fiber.Ctx) error {

	id := c.Params("id")

	data, err := service.GetSiswaByID(id)

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

	// VALIDASI
	if siswa.Nama == "" ||
		siswa.Kelas == "" ||
		siswa.Alamat == "" ||
		siswa.Email == "" {

		return helpers.ErrorResponse(c, 400, "Semua field (nama, kelas, alamat, email) wajib diisi")
	}

	if !isValidEmail(siswa.Email) {
		return helpers.ErrorResponse(c, 400, "Format email tidak valid")
	}

	if service.CheckEmailExists(siswa.Email, "") {
		return helpers.ErrorResponse(c, 400, "Email sudah terdaftar, gunakan email lain")
	}

	data, err := service.CreateSiswa(siswa)

	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal menambahkan siswa")
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Siswa berhasil ditambahkan",
		"data":    data,
	})
}

func UpdateSiswa(c *fiber.Ctx) error {

	id := c.Params("id")

	var siswa model.Siswa

	err := c.BodyParser(&siswa)

	if err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	// VALIDASI UPDATE
	if siswa.Nama == "" ||
		siswa.Kelas == "" ||
		siswa.Alamat == "" ||
		siswa.Email == "" {

		return helpers.ErrorResponse(c, 400, "Semua field (nama, kelas, alamat, email) wajib diisi")
	}

	if !isValidEmail(siswa.Email) {
		return helpers.ErrorResponse(c, 400, "Format email tidak valid")
	}

	if service.CheckEmailExists(siswa.Email, id) {
		return helpers.ErrorResponse(c, 400, "Email sudah terdaftar pada siswa lain")
	}

	data, err := service.UpdateSiswa(id, siswa)

	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal update siswa")
	}

	return helpers.SuccessResponse(c, "Siswa berhasil diupdate", data)
}

func DeleteSiswa(c *fiber.Ctx) error {

	id := c.Params("id")

	err := service.DeleteSiswa(id)

	if err != nil {
		return helpers.ErrorResponse(c, 500, "Gagal menghapus siswa")
	}

	return helpers.SuccessResponse(c, "Siswa berhasil dihapus", nil)
}