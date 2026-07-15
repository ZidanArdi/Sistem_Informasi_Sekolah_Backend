package handler

import (
	"errors"

	"backend/config"
	"backend/helpers"
	authModel "backend/modules/auth/model"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserResponse struct {
	ID         uint   `json:"id" example:"1"`
	Name       string `json:"name" example:"Ahmad Guru"`
	Role       string `json:"role" example:"guru"`
	Identifier string `json:"identifier" example:"GR001"`
	IsActive   bool   `json:"is_active" example:"true"`
}

type UpdateStatusInput struct {
	IsActive *bool `json:"is_active" example:"true"`
}

// GetUsers godoc
// @Summary Ambil semua data user account
// @Description Mengambil daftar seluruh user account (Admin, Guru, Siswa) beserta nama asli dan identifier (NIP/NIS/Email). Hanya dapat diakses oleh Admin.
// @Tags User Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} helpers.SwaggerSuccessResponse{data=[]UserResponse}
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 403 {object} helpers.SwaggerError403Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/users/ [get]
func GetUsers(c *fiber.Ctx) error {
	var users []authModel.User
	if err := config.DB.Order("id desc").Find(&users).Error; err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data user: "+err.Error())
	}

	type SiswaInfo struct {
		Nama string
		NIS  string
	}
	var siswas []struct {
		UserID uint
		Nama   string
		NIS    string
	}
	config.DB.Table("siswas").Select("user_id, nama, nis").Find(&siswas)
	siswaMap := make(map[uint]SiswaInfo)
	for _, s := range siswas {
		siswaMap[s.UserID] = SiswaInfo{Nama: s.Nama, NIS: s.NIS}
	}

	var gurus []struct {
		UserID uint
		Nama   string
		NIP    string `gorm:"column:nip"`
	}
	config.DB.Table("gurus").Select("user_id, nama, nip").Find(&gurus)
	type GuruInfo struct {
		Nama string
		NIP  string
	}
	guruMap := make(map[uint]GuruInfo)
	for _, g := range gurus {
		guruMap[g.UserID] = GuruInfo{Nama: g.Nama, NIP: g.NIP}
	}

	response := make([]UserResponse, 0)
	for _, u := range users {
		name := u.Nama
		identifier := ""
		if u.Email != nil {
			identifier = *u.Email
		}

		if u.Role == "siswa" {
			if info, exists := siswaMap[u.ID]; exists {
				name = info.Nama
				identifier = info.NIS
			}
		} else if u.Role == "guru" {
			if info, exists := guruMap[u.ID]; exists {
				name = info.Nama
				identifier = info.NIP
			}
		}

		response = append(response, UserResponse{
			ID:         u.ID,
			Name:       name,
			Role:       u.Role,
			Identifier: identifier,
			IsActive:   u.IsActive,
		})
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil data user", response)
}

// UpdateUserStatus godoc
// @Summary Perbarui status keaktifan user
// @Description Mengaktifkan atau menonaktifkan user account berdasarkan ID. Hanya dapat diakses oleh Admin. Administrator tidak dapat dinonaktifkan.
// @Tags User Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID User"
// @Param request body UpdateStatusInput true "Payload status aktif"
// @Success 200 {object} helpers.SwaggerSuccessResponse{data=authModel.User}
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 403 {object} helpers.SwaggerError403Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/users/{id}/status [put]
func UpdateUserStatus(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	var input UpdateStatusInput
	if err := c.BodyParser(&input); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Format payload tidak valid")
	}

	if input.IsActive == nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "is_active wajib diisi")
	}

	var user authModel.User
	if err := config.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.ErrorResponse(c, fiber.StatusNotFound, "User tidak ditemukan")
		}
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if user.Role == "admin" {
		return helpers.ErrorResponse(c, fiber.StatusForbidden, "ERR_PERMISSION_DENIED: Status administrator tidak dapat diubah")
	}

	user.IsActive = *input.IsActive
	if err := config.DB.Save(&user).Error; err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengupdate status user: "+err.Error())
	}

	return helpers.SuccessResponse(c, "Status user berhasil diperbarui", user)
}

// ResetUserPassword godoc
// @Summary Reset password user ke password default
// @Description Mereset password user berdasarkan ID ke password default berdasarkan role (e.g. Guru123! untuk guru, Siswa123! untuk siswa). Hanya dapat diakses oleh Admin.
// @Tags User Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID User"
// @Success 200 {object} helpers.SwaggerSuccessResponse
// @Failure 400 {object} helpers.SwaggerError400Response
// @Failure 401 {object} helpers.SwaggerError401Response
// @Failure 403 {object} helpers.SwaggerError403Response
// @Failure 404 {object} helpers.SwaggerError404Response
// @Failure 500 {object} helpers.SwaggerError500Response
// @Router /api/users/{id}/reset-password [put]
func ResetUserPassword(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	var user authModel.User
	if err := config.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.ErrorResponse(c, fiber.StatusNotFound, "User tidak ditemukan")
		}
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if user.Role == "admin" {
		return helpers.ErrorResponse(c, fiber.StatusForbidden, "ERR_PERMISSION_DENIED: Password administrator tidak dapat direset")
	}

	var plainPassword string
	if user.Role == "admin" {
		plainPassword = "Admin123!"
	} else if user.Role == "guru" {
		plainPassword = "Guru123!"
	} else {
		plainPassword = "Siswa123!"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal melakukan hashing password: "+err.Error())
	}

	user.Password = string(hashedPassword)
	user.IsFirstLogin = true // For student to change password again upon reset if they are student
	if err := config.DB.Save(&user).Error; err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal menyimpan password baru: "+err.Error())
	}

	return helpers.SuccessResponse(c, "Password user berhasil direset", fiber.Map{
		"password": plainPassword,
	})
}
