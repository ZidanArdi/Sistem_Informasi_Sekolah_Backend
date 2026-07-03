package service

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"backend/config"
	authModel "backend/modules/auth/model"
	"backend/modules/siswa/model"
	"backend/modules/siswa/repository"

	"golang.org/x/crypto/bcrypt"
)

func GetAllSiswa(search string, kelasID string) ([]model.Siswa, error) {
	return repository.GetAllSiswa(search, kelasID)
}

func GetSiswaByID(id uint) (model.Siswa, error) {
	return repository.GetSiswaByID(id)
}

func GenerateRandomPassword() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(90000) + 10000 // Generates a number between 10000 and 99999
	return fmt.Sprintf("STD-%d", num)
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
	nis, err := repository.GenerateNISWithTx(tx)
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

	// 3. Generate random temporary password
	plainPassword := GenerateRandomPassword()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return model.Siswa{}, "", err
	}

	// 4. Create auth user record (Email is NULL for student)
	user := authModel.User{
		Nama:         data.Nama,
		Email:        nil,
		Password:     string(hashedPassword),
		Role:         "siswa",
		IsFirstLogin: true, // Required for force change password
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return model.Siswa{}, "", errors.New("gagal membuat akun user untuk siswa: " + err.Error())
	}

	// 5. Link User ID to Siswa record and save Siswa
	data.UserID = user.ID
	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return model.Siswa{}, "", errors.New("gagal menyimpan data siswa: " + err.Error())
	}

	// Commit Transaction
	if err := tx.Commit().Error; err != nil {
		return model.Siswa{}, "", err
	}

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
		strings.TrimSpace(data.TempatLahir) == "" ||
		strings.TrimSpace(data.TanggalLahir) == "" ||
		strings.TrimSpace(data.Alamat) == "" ||
		data.KelasID == 0 {
		return errors.New("nama, jenis_kelamin, tempat_lahir, tanggal_lahir, alamat, dan kelas_id wajib diisi")
	}
	return nil
}
