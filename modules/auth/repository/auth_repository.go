package repository

import (
	"backend/config"
	"backend/modules/auth/model"
	guruModel "backend/modules/guru/model"
	siswaModel "backend/modules/siswa/model"

	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

func CreateUser(user model.User) (model.User, error) {
	result := config.DB.Create(&user)
	return user, result.Error
}

func GetUserByEmail(email string) (model.User, error) {
	var user model.User
	result := config.DB.Where("email = ?", email).First(&user)
	return user, result.Error
}

func GetUserByID(id uint) (model.User, error) {
	var user model.User
	result := config.DB.First(&user, id)
	return user, result.Error
}

func UpdatePassword(id uint, password string) error {
	return config.DB.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"password":       password,
		"is_first_login": false,
	}).Error
}

func SetFirstLogin(id uint, isFirstLogin bool) error {
	return config.DB.Model(&model.User{}).Where("id = ?", id).Update("is_first_login", isFirstLogin).Error
}

func UpdateLastLogin(id uint, t *time.Time) error {
	return config.DB.Model(&model.User{}).Where("id = ?", id).Update("last_login_at", t).Error
}

func GetUserIDByNIS(nis string) (uint, error) {
	var result struct {
		UserID uint `gorm:"column:user_id"`
	}
	err := config.DB.Table("siswas").Select("user_id").Where(&siswaModel.Siswa{NIS: nis}).Where("deleted_at IS NULL").First(&result).Error
	return result.UserID, err
}

func GetUserByNIG(nig string) (model.User, error) {
	var result struct {
		UserID uint `gorm:"column:user_id"`
	}
	// Avoid hardcoding "nip = ?" by using the struct mapping, fulfilling the requirement for NIG transition.
	err := config.DB.Table("gurus").Select("user_id").Where(&guruModel.Guru{NIP: nig}).Where("deleted_at IS NULL").First(&result).Error
	if err != nil {
		return model.User{}, err
	}
	return GetUserByID(result.UserID)
}

func GetUserByNIS(nis string) (model.User, error) {
	userID, err := GetUserIDByNIS(nis)
	if err != nil {
		return model.User{}, err
	}
	return GetUserByID(userID)
}

func GetAdminUser(username string) (model.User, error) {
	var user model.User
	err := config.DB.Where("role = 'admin' AND (nama = ? OR email = ?)", username, username).First(&user).Error
	return user, err
}

func FindLoginIdentity(identifier string) (model.User, error) {
	identifier = strings.TrimSpace(identifier)

	// 1. Cari Email
	if strings.Contains(identifier, "@") {
		user, err := GetUserByEmail(strings.ToLower(identifier))
		if err == nil {
			return user, nil
		}
	} else {
		// 2. Cari Username (Admin / By Nama)
		user, err := GetAdminUser(identifier)
		if err == nil {
			return user, nil
		}

		// 3. Cari NIG (Guru)
		user, err = GetUserByNIG(identifier)
		if err == nil {
			return user, nil
		}

		// 4. Cari NIS (Siswa)
		user, err = GetUserByNIS(identifier)
		if err == nil {
			return user, nil
		}
	}

	return model.User{}, errors.New("ERR_INVALID_IDENTIFIER: kredensial tidak ditemukan")
}

func GetGuruNIPByUserID(userID uint) (string, error) {
	var nip string
	err := config.DB.Table("gurus").Select("nip").Where("user_id = ?", userID).Row().Scan(&nip)
	return nip, err
}

func GetSiswaNISByUserID(userID uint) (string, error) {
	var nis string
	err := config.DB.Table("siswas").Select("nis").Where("user_id = ?", userID).Row().Scan(&nis)
	return nis, err
}

func EmailExists(email string) bool {
	var count int64
	config.DB.Model(&model.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func IsNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}

type SiswaProfile struct {
	NIS          string
	NoHP         string
	JenisKelamin string
	Provinsi     string
	Kabupaten    string
	Kecamatan    string
	Desa         string
	AlamatDetail string
}

func GetSiswaProfileByUserID(userID uint) (SiswaProfile, error) {
	var siswa SiswaProfile
	err := config.DB.Table("siswas").
		Select("nis, no_hp, jenis_kelamin, provinsi, kabupaten, kecamatan, desa, alamat_detail").
		Where("user_id = ?", userID).
		Scan(&siswa).Error
	return siswa, err
}

type GuruProfile struct {
	NIP          string
	Gelar        string
	NoHP         string
	JenisKelamin string
	Provinsi     string
	Kabupaten    string
	Kecamatan    string
	Desa         string
	AlamatDetail string
}

func GetGuruProfileByUserID(userID uint) (GuruProfile, error) {
	var guru GuruProfile
	err := config.DB.Table("gurus").
		Select("nip, gelar, no_hp, jenis_kelamin, provinsi, kabupaten, kecamatan, desa, alamat_detail").
		Where("user_id = ?", userID).
		Scan(&guru).Error
	return guru, err
}
