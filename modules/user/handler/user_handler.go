package handler

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"backend/config"
	"backend/helpers"
	authModel "backend/modules/auth/model"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	Identifier string `json:"identifier"`
	IsActive   bool   `json:"is_active"`
}

type UpdateStatusInput struct {
	IsActive *bool `json:"is_active" xml:"is_active" form:"is_active"`
}

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
	}
	config.DB.Table("gurus").Select("user_id, nama").Find(&gurus)
	guruMap := make(map[uint]string)
	for _, g := range gurus {
		guruMap[g.UserID] = g.Nama
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
			if gName, exists := guruMap[u.ID]; exists {
				name = gName
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

	user.IsActive = *input.IsActive
	if err := config.DB.Save(&user).Error; err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengupdate status user: "+err.Error())
	}

	return helpers.SuccessResponse(c, "Status user berhasil diperbarui", user)
}

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

	// Generate random password like RST-XXXXX
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(90000) + 10000
	plainPassword := fmt.Sprintf("RST-%d", num)

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
