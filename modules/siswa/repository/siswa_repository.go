package repository

import (
	"backend/config"
	"backend/modules/siswa/model"
	"strconv"

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

	siswa.NIS = data.NIS
	siswa.Nama = data.Nama
	siswa.JenisKelamin = data.JenisKelamin
	siswa.TempatLahir = data.TempatLahir
	siswa.TanggalLahir = data.TanggalLahir
	siswa.Alamat = data.Alamat
	siswa.Email = data.Email
	siswa.KelasID = data.KelasID

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
