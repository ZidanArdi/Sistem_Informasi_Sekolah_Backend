package repository

import (
	"backend/config"
	"backend/modules/siswa/model"

	"gorm.io/gorm"
)

func GetAllSiswa(search string) ([]model.Siswa, error) {

	var siswa []model.Siswa

	query := config.DB

	if search != "" {
		query = query.Where("nama ILIKE ?", "%"+search+"%")
	}

	result := query.Find(&siswa)

	return siswa, result.Error
}

func GetSiswaByID(id string) (model.Siswa, error) {

	var siswa model.Siswa

	result := config.DB.Where("nis = ?", id).First(&siswa)

	return siswa, result.Error
}

func CreateSiswa(data model.Siswa) (model.Siswa, error) {

	result := config.DB.Create(&data)

	return data, result.Error
}

func UpdateSiswa(id string, data model.Siswa) (model.Siswa, error) {

	var siswa model.Siswa

	err := config.DB.Where("nis = ?", id).First(&siswa).Error

	if err != nil {
		return siswa, err
	}

	siswa.Nama = data.Nama
	siswa.Kelas = data.Kelas
	siswa.Alamat = data.Alamat
	siswa.Email = data.Email

	config.DB.Save(&siswa)

	return siswa, nil
}

func DeleteSiswa(id string) error {

	var siswa model.Siswa

	err := config.DB.Where("nis = ?", id).First(&siswa).Error

	if err != nil {
		return err
	}

	return config.DB.Delete(&siswa).Error
}

func IsNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}

func CheckEmailExists(email string, excludeNIS string) bool {
	var count int64
	query := config.DB.Model(&model.Siswa{}).Where("email = ?", email)
	if excludeNIS != "" {
		query = query.Where("nis != ?", excludeNIS)
	}
	query.Count(&count)
	return count > 0
}