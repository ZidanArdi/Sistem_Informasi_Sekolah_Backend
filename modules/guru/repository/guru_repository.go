package repository

import (
	"backend/config"
	"backend/modules/guru/model"

	"gorm.io/gorm"
)

func GetAllGuru(search string) ([]model.Guru, error) {
	var guru []model.Guru

	query := config.DB
	if search != "" {
		query = query.Where("nama ILIKE ? OR nip ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	result := query.Find(&guru)
	return guru, result.Error
}

func GetGuruByID(id uint) (model.Guru, error) {
	var guru model.Guru
	result := config.DB.First(&guru, id)
	return guru, result.Error
}

func CreateGuru(data model.Guru) (model.Guru, error) {
	result := config.DB.Create(&data)
	return data, result.Error
}

func UpdateGuru(id uint, data model.Guru) (model.Guru, error) {
	var guru model.Guru

	if err := config.DB.First(&guru, id).Error; err != nil {
		return guru, err
	}

	guru.NIP = data.NIP
	guru.Nama = data.Nama
	guru.JenisKelamin = data.JenisKelamin
	guru.Email = data.Email
	guru.NoHP = data.NoHP
	guru.Alamat = data.Alamat

	err := config.DB.Save(&guru).Error
	return guru, err
}

func DeleteGuru(id uint) error {
	var guru model.Guru

	if err := config.DB.First(&guru, id).Error; err != nil {
		return err
	}

	return config.DB.Delete(&guru).Error
}

func CheckNIPExists(nip string, excludeID uint) bool {
	var count int64
	query := config.DB.Model(&model.Guru{}).Where("nip = ?", nip)
	if excludeID != 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count > 0
}

func IsNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}
