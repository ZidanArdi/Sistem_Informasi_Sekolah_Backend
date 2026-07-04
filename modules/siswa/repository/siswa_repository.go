package repository

import (
	"backend/config"
	"backend/modules/siswa/model"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func GetAllSiswa(search string, kelasID string) ([]model.Siswa, error) {

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

	result := query.Find(&siswa)

	return siswa, result.Error
}

func GetSiswaByID(id uint) (model.Siswa, error) {

	var siswa model.Siswa

	result := config.DB.Preload("Kelas").Preload("Kelas.WaliKelas").First(&siswa, id)

	return siswa, result.Error
}

func CreateSiswa(data model.Siswa) (model.Siswa, error) {

	result := config.DB.Create(&data)

	return data, result.Error
}

func UpdateSiswa(id uint, data model.Siswa) (model.Siswa, error) {

	var siswa model.Siswa

	err := config.DB.First(&siswa, id).Error

	if err != nil {
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

	err := config.DB.First(&siswa, id).Error

	if err != nil {
		return err
	}

	return config.DB.Delete(&siswa).Error
}

func IsNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}

func CheckNISExists(nis string, excludeID uint) bool {
	var count int64
	query := config.DB.Model(&model.Siswa{}).Where("nis = ?", nis)
	if excludeID != 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count > 0
}

func GenerateNISWithTx(tx *gorm.DB) (string, error) {
	year := time.Now().Year()
	prefix := fmt.Sprintf("%d", year)

	var latestSiswa model.Siswa
	// Lock the row for update to prevent concurrent race condition duplicates
	err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("nis LIKE ?", prefix+"%").
		Order("nis desc").
		First(&latestSiswa).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return prefix + "0001", nil
		}
		return "", err
	}

	var seq int
	_, err = fmt.Sscanf(latestSiswa.NIS, prefix+"%d", &seq)
	if err != nil {
		return prefix + "0001", nil
	}

	return fmt.Sprintf("%s%04d", prefix, seq+1), nil
}
