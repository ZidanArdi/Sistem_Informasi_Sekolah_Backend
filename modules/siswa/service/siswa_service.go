package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"backend/config"
	authModel "backend/modules/auth/model"
	kelasModel "backend/modules/kelas/model"
	"backend/modules/siswa/model"
	"backend/modules/siswa/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetAllSiswa(search string, kelasID string, guruID string) ([]model.Siswa, error) {
	return repository.GetAllSiswa(search, kelasID, guruID)
}

func GetSiswaByID(id uint) (model.Siswa, error) {
	return repository.GetSiswaByID(id)
}

func GenerateNISWithTx(tx *gorm.DB) (string, error) {
	var latestSiswa model.Siswa
	err := tx.Set("gorm:query_option", "FOR UPDATE").
		Unscoped().
		Order("nis desc").
		First(&latestSiswa).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "20260001", nil
		}
		return "", fmt.Errorf("ERR_GENERATE_NIS: %v", err)
	}

	seq, err := strconv.ParseInt(latestSiswa.NIS, 10, 64)
	if err != nil {
		return "20260001", nil
	}

	return fmt.Sprintf("%d", seq+1), nil
}

func CreateSiswa(data model.Siswa) (model.Siswa, string, error) {
	// Start Database Transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Auto generate NIS using safe database lock transaction
	nis, err := GenerateNISWithTx(tx)
	if err != nil {
		tx.Rollback()
		return model.Siswa{}, "", err
	}
	data.NIS = nis

	// 2. Validate siswa input fields
	if err := validateSiswa(data); err != nil {
		tx.Rollback()
		return model.Siswa{}, "", err
	}

	// Generate student email dynamically based on NIS suffix
	var seq int
	if len(data.NIS) >= 4 {
		seqStr := data.NIS[len(data.NIS)-4:]
		seq, _ = strconv.Atoi(seqStr)
	}
	if seq == 0 {
		seq = 1
	}
	email := fmt.Sprintf("siswa%03d@sekolah.com", seq)

	// Check if NIS already exists
	var exists int64
	if err := tx.Model(&model.Siswa{}).Where("nis = ?", data.NIS).Count(&exists).Error; err != nil {
		tx.Rollback()
		return model.Siswa{}, "", fmt.Errorf("ERR_GENERATE_NIS: Gagal memeriksa keunikan NIS: %v", err)
	}
	if exists > 0 {
		tx.Rollback()
		return model.Siswa{}, "", errors.New("ERR_IDENTIFIER_DUPLICATE: NIS sudah terdaftar")
	}

	// Check if Email already exists
	var emailCount int64
	if err := tx.Model(&authModel.User{}).Where("email = ?", email).Count(&emailCount).Error; err != nil {
		tx.Rollback()
		return model.Siswa{}, "", fmt.Errorf("ERR_GENERATE_NIS: Gagal memeriksa keunikan email: %v", err)
	}
	if emailCount > 0 {
		tx.Rollback()
		return model.Siswa{}, "", errors.New("ERR_EMAIL_DUPLICATE: Email sudah terdaftar")
	}

	// 3. Set default password
	plainPassword := "Siswa123!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return model.Siswa{}, "", err
	}

	// 4. Create auth user record
	user := authModel.User{
		Nama:         data.Nama,
		Email:        &email,
		Password:     string(hashedPassword),
		Role:         "siswa",
		IsFirstLogin: true,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return model.Siswa{}, "", errors.New("ERR_EMAIL_DUPLICATE: gagal membuat akun user untuk siswa: " + err.Error())
	}

	// 5. Link User ID to Siswa record and save Siswa
	data.UserID = user.ID
	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return model.Siswa{}, "", errors.New("ERR_GENERATE_NIS: gagal menyimpan data siswa: " + err.Error())
	}

	// Commit Transaction
	if err := tx.Commit().Error; err != nil {
		return model.Siswa{}, "", err
	}

	// LOG Student Created Info
	log.Printf("[INFO] Action: Student Created | Timestamp: %s | User: %s | Identifier: %s | Email: %s", time.Now().Format(time.RFC3339), data.Nama, data.NIS, email)

	// Preload class info for response
	config.DB.Preload("Kelas").First(&data, data.ID)

	return data, plainPassword, nil
}

func UpdateSiswa(id uint, data model.Siswa) (model.Siswa, error) {
	if err := validateSiswa(data); err != nil {
		return model.Siswa{}, err
	}
	return repository.UpdateSiswa(id, data)
}

func DeleteSiswa(id uint) error {
	return repository.DeleteSiswa(id)
}

func validateSiswa(data model.Siswa) error {
	if strings.TrimSpace(data.Nama) == "" ||
		strings.TrimSpace(data.JenisKelamin) == "" ||
		strings.TrimSpace(data.TanggalLahir) == "" ||
		strings.TrimSpace(data.Provinsi) == "" ||
		strings.TrimSpace(data.Kabupaten) == "" ||
		strings.TrimSpace(data.Kecamatan) == "" ||
		strings.TrimSpace(data.Desa) == "" ||
		strings.TrimSpace(data.AlamatDetail) == "" ||
		data.KelasID == 0 {
		return errors.New("nama, jenis_kelamin, tanggal_lahir, provinsi, kabupaten, kecamatan, desa, alamat_detail, dan kelas_id wajib diisi")
	}

	var kelas kelasModel.Kelas
	if err := config.DB.First(&kelas, data.KelasID).Error; err != nil {
		return errors.New("kelas tidak ditemukan")
	}

	var count int64
	query := config.DB.Model(&model.Siswa{}).Where("kelas_id = ?", data.KelasID)
	if data.ID != 0 {
		query = query.Where("id != ?", data.ID)
	}
	if err := query.Count(&count).Error; err != nil {
		return err
	}

	if int(count) >= kelas.Kapasitas {
		return errors.New("Kelas sudah mencapai kapasitas maksimum.")
	}

	return nil
}
