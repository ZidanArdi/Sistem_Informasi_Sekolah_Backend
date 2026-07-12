package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"backend/config"
	"backend/helpers"
	guruModel "backend/modules/guru/model"
	siswaModel "backend/modules/siswa/model"

	"github.com/gofiber/fiber/v2"
)

type ProfileUpdateInput struct {
	NoHP         string `json:"no_hp"`
	Provinsi     string `json:"provinsi"`
	Kabupaten    string `json:"kabupaten"`
	Kecamatan    string `json:"kecamatan"`
	Desa         string `json:"desa"`
	AlamatDetail string `json:"alamat_detail"`
}

// GetProfile returns the profile of the currently logged-in user
func GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	role, okRole := c.Locals("role").(string)
	if !ok || !okRole {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "User tidak valid")
	}

	if role == "guru" {
		var guru guruModel.Guru
		// Preload User details like email
		if err := config.DB.Preload("User").Where("user_id = ?", userID).First(&guru).Error; err != nil {
			return helpers.ErrorResponse(c, fiber.StatusNotFound, "Data profil guru tidak ditemukan")
		}

		// Fetch mapel names taught by this guru
		var mapelNames []string
		config.DB.Table("guru_mapels").
			Select("mapels.nama").
			Joins("join mapels on mapels.id = guru_mapels.mapel_id").
			Where("guru_mapels.guru_id = ?", guru.ID).
			Pluck("nama", &mapelNames)

		emailStr := ""
		var lastLogin *time.Time
		if guru.User != nil {
			if guru.User.Email != nil {
				emailStr = *guru.User.Email
			}
			lastLogin = guru.User.LastLoginAt
		}

		return helpers.SuccessResponse(c, "Berhasil mengambil profil guru", fiber.Map{
			"role":            "guru",
			"id":              guru.ID,
			"nip":             guru.NIP,
			"nama":            guru.Nama,
			"gelar":           guru.Gelar,
			"jenis_kelamin":   guru.JenisKelamin,
			"no_hp":           guru.NoHP,
			"provinsi":        guru.Provinsi,
			"kabupaten":       guru.Kabupaten,
			"kecamatan":       guru.Kecamatan,
			"desa":            guru.Desa,
			"alamat_detail":   guru.AlamatDetail,
			"photo_url":       guru.PhotoURL,
			"email":           emailStr,
			"subjects_taught": mapelNames,
			"last_login_at":   lastLogin,
		})
	} else if role == "siswa" {
		var siswa siswaModel.Siswa
		if err := config.DB.Preload("Kelas").Where("user_id = ?", userID).First(&siswa).Error; err != nil {
			return helpers.ErrorResponse(c, fiber.StatusNotFound, "Data profil siswa tidak ditemukan")
		}

		kelasNama := "-"
		if siswa.Kelas.NamaKelas != "" {
			kelasNama = siswa.Kelas.NamaKelas
		}

		var lastLogin *time.Time
		config.DB.Table("users").Select("last_login_at").Where("id = ?", userID).Scan(&lastLogin)

		return helpers.SuccessResponse(c, "Berhasil mengambil profil siswa", fiber.Map{
			"role":          "siswa",
			"id":            siswa.ID,
			"nis":           siswa.NIS,
			"nama":          siswa.Nama,
			"jenis_kelamin": siswa.JenisKelamin,
			"tanggal_lahir": siswa.TanggalLahir,
			"no_hp":         siswa.NoHP,
			"provinsi":      siswa.Provinsi,
			"kabupaten":     siswa.Kabupaten,
			"kecamatan":     siswa.Kecamatan,
			"desa":          siswa.Desa,
			"alamat_detail": siswa.AlamatDetail,
			"photo_url":     siswa.PhotoURL,
			"kelas":         kelasNama,
			"last_login_at": lastLogin,
		})
	}

	// Fallback for Admin
	var user struct {
		Nama        string
		Email       string
		LastLoginAt *time.Time
	}
	config.DB.Table("users").Select("nama, email, last_login_at").Where("id = ?", userID).Scan(&user)

	return helpers.SuccessResponse(c, "Berhasil mengambil profil admin", fiber.Map{
		"role":          role,
		"nama":          user.Nama,
		"email":         user.Email,
		"last_login_at": user.LastLoginAt,
	})
}

// UpdateProfile updates the profile of the currently logged-in user (Guru/Siswa)
func UpdateProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	role, okRole := c.Locals("role").(string)
	if !ok || !okRole {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "User tidak valid")
	}

	var input ProfileUpdateInput
	if err := c.BodyParser(&input); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Input tidak valid")
	}

	if role == "guru" {
		var guru guruModel.Guru
		if err := config.DB.Where("user_id = ?", userID).First(&guru).Error; err != nil {
			return helpers.ErrorResponse(c, fiber.StatusNotFound, "Data guru tidak ditemukan")
		}

		guru.NoHP = strings.TrimSpace(input.NoHP)
		guru.Provinsi = strings.TrimSpace(input.Provinsi)
		guru.Kabupaten = strings.TrimSpace(input.Kabupaten)
		guru.Kecamatan = strings.TrimSpace(input.Kecamatan)
		guru.Desa = strings.TrimSpace(input.Desa)
		guru.AlamatDetail = strings.TrimSpace(input.AlamatDetail)

		if err := config.DB.Save(&guru).Error; err != nil {
			return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengupdate profil")
		}

		return helpers.SuccessResponse(c, "Profil berhasil diupdate", guru)
	} else if role == "siswa" {
		var siswa siswaModel.Siswa
		if err := config.DB.Where("user_id = ?", userID).First(&siswa).Error; err != nil {
			return helpers.ErrorResponse(c, fiber.StatusNotFound, "Data siswa tidak ditemukan")
		}

		siswa.NoHP = strings.TrimSpace(input.NoHP)
		siswa.Provinsi = strings.TrimSpace(input.Provinsi)
		siswa.Kabupaten = strings.TrimSpace(input.Kabupaten)
		siswa.Kecamatan = strings.TrimSpace(input.Kecamatan)
		siswa.Desa = strings.TrimSpace(input.Desa)
		siswa.AlamatDetail = strings.TrimSpace(input.AlamatDetail)

		if err := config.DB.Save(&siswa).Error; err != nil {
			return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengupdate profil")
		}

		return helpers.SuccessResponse(c, "Profil berhasil diupdate", siswa)
	}

	return helpers.ErrorResponse(c, fiber.StatusForbidden, "Role tidak diperbolehkan mengupdate profil")
}

// UploadProfilePhoto handles profile photo uploads and validations
func UploadProfilePhoto(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	role, okRole := c.Locals("role").(string)
	if !ok || !okRole {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "User tidak valid")
	}

	// 1. Retrieve the file
	file, err := c.FormFile("photo")
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "File photo tidak ditemukan")
	}

	// 2. Validate Size (Max 2MB)
	if file.Size > 2*1024*1024 {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Format file tidak didukung atau ukuran file terlalu besar.")
	}

	// 3. Validate Format/Extension (jpg, jpeg, png)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Format file tidak didukung atau ukuran file terlalu besar.")
	}

	// 4. Build filename based on strategy: role_userid_timestamp.ext
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%s_%d_%d%s", role, userID, timestamp, ext)

	// 5. Ensure storage directory exists
	uploadDir := "./uploads/profile"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal menyiapkan direktori penyimpanan")
	}

	// 6. Save the file
	savePath := filepath.Join(uploadDir, filename)
	if err := c.SaveFile(file, savePath); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal menyimpan file photo")
	}

	photoURL := "/uploads/profile/" + filename

	// 7. Update database based on role
	if role == "guru" {
		var guru guruModel.Guru
		if err := config.DB.Where("user_id = ?", userID).First(&guru).Error; err != nil {
			return helpers.ErrorResponse(c, fiber.StatusNotFound, "Data guru tidak ditemukan")
		}
		guru.PhotoURL = photoURL
		config.DB.Save(&guru)
	} else if role == "siswa" {
		var siswa siswaModel.Siswa
		if err := config.DB.Where("user_id = ?", userID).First(&siswa).Error; err != nil {
			return helpers.ErrorResponse(c, fiber.StatusNotFound, "Data siswa tidak ditemukan")
		}
		siswa.PhotoURL = photoURL
		config.DB.Save(&siswa)
	} else {
		return helpers.ErrorResponse(c, fiber.StatusForbidden, "Role tidak diperbolehkan mengupload photo")
	}

	return helpers.SuccessResponse(c, "Photo profil berhasil diupload", fiber.Map{
		"photo_url": photoURL,
	})
}
