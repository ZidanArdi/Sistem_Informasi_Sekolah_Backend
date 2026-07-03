package repository

import (
	"backend/config"
	"backend/modules/guru/model"
	"fmt"

	"gorm.io/gorm"
)

func GetAllGuru(search string) ([]model.Guru, error) {
	var guru []model.Guru

	query := config.DB.Preload("User")
	if search != "" {
		query = query.Where("nama ILIKE ? OR nip ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	result := query.Find(&guru)
	return guru, result.Error
}

func GetGuruByID(id uint) (model.Guru, error) {
	var guru model.Guru
	result := config.DB.Preload("User").First(&guru, id)
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

	guru.Nama = data.Nama
	guru.JenisKelamin = data.JenisKelamin
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

func GenerateNIPWithTx(tx *gorm.DB) (string, error) {
	prefix := "GURU"

	var latestGuru model.Guru
	err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("nip LIKE ?", prefix+"%").
		Order("nip desc").
		First(&latestGuru).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return prefix + "0001", nil
		}
		return "", err
	}

	var seq int
	_, err = fmt.Sscanf(latestGuru.NIP, prefix+"%04d", &seq)
	if err != nil {
		return prefix + "0001", nil
	}

	return fmt.Sprintf("%s%04d", prefix, seq+1), nil
}
