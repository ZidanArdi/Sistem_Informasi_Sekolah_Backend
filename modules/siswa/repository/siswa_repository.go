package repository

import (
	"fmt"
	"strconv"

	"backend/config"
	authModel "backend/modules/auth/model"
	kelasModel "backend/modules/kelas/model"
	"backend/modules/siswa/model"

	"gorm.io/gorm"
)

func BeginTransaction() *gorm.DB {
	return config.DB.Begin()
}

func GetAllSiswa(search string, kelasID string, guruID string) ([]model.Siswa, error) {
	var siswa []model.Siswa
	query := config.DB.Preload("Kelas").Preload("Kelas.WaliKelas")
	if search != "" {
		query = query.Where("nama ILIKE ? OR nis ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if kelasID != "" {
		if parsedKelasID, err := strconv.Atoi(kelasID); err == nil {
			query = query.Where("kelas_id = ?", parsedKelasID)
		}
	}
	if guruID != "" {
		if parsedGuruID, err := strconv.Atoi(guruID); err == nil {
			query = query.Where("kelas_id IN (SELECT DISTINCT kelas_id FROM jadwals WHERE guru_id = ? AND deleted_at IS NULL)", parsedGuruID)
		}
	}
	result := query.Find(&siswa)
	return siswa, result.Error
}

func GetSiswaByID(id uint) (model.Siswa, error) {
	var siswa model.Siswa
	result := config.DB.Preload("Kelas").Preload("Kelas.WaliKelas").First(&siswa, id)
	return siswa, result.Error
}

func CreateSiswaWithTx(tx *gorm.DB, data model.Siswa) error {
	return tx.Create(&data).Error
}

func CreateUserWithTx(tx *gorm.DB, user *authModel.User) error {
	return tx.Create(user).Error
}

func UpdateSiswa(id uint, data model.Siswa) (model.Siswa, error) {
	var siswa model.Siswa
	if err := config.DB.First(&siswa, id).Error; err != nil {
		return siswa, err
	}
	siswa.Nama = data.Nama
	siswa.JenisKelamin = data.JenisKelamin
	siswa.TanggalLahir = data.TanggalLahir
	siswa.NoHP = data.NoHP
	siswa.KelasID = data.KelasID
	siswa.UserID = data.UserID
	siswa.Provinsi = data.Provinsi
	siswa.Kabupaten = data.Kabupaten
	siswa.Kecamatan = data.Kecamatan
	siswa.Desa = data.Desa
	siswa.AlamatDetail = data.AlamatDetail
	if err := config.DB.Save(&siswa).Error; err != nil {
		return siswa, err
	}
	config.DB.Preload("Kelas").Preload("Kelas.WaliKelas").First(&siswa, siswa.ID)
	return siswa, nil
}

func DeleteSiswa(id uint) error {
	var siswa model.Siswa
	if err := config.DB.First(&siswa, id).Error; err != nil {
		return err
	}
	return config.DB.Delete(&siswa).Error
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

func CountSiswaByNISWithTx(tx *gorm.DB, nis string) (int64, error) {
	var count int64
	err := tx.Model(&model.Siswa{}).Where("nis = ?", nis).Count(&count).Error
	return count, err
}

func CountUserByEmailWithTx(tx *gorm.DB, email string) (int64, error) {
	var count int64
	err := tx.Model(&authModel.User{}).Where("email = ?", email).Count(&count).Error
	return count, err
}

func GetKelasByID(id uint) (kelasModel.Kelas, error) {
	var kelas kelasModel.Kelas
	err := config.DB.First(&kelas, id).Error
	return kelas, err
}

func CountSiswaByKelasExcludeID(kelasID uint, excludeID uint) (int64, error) {
	var count int64
	query := config.DB.Model(&model.Siswa{}).Where("kelas_id = ?", kelasID)
	if excludeID != 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count, err
}

func LoadSiswaRelations(siswa *model.Siswa) {
	config.DB.Preload("Kelas").First(siswa, siswa.ID)
}

func IsNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}
