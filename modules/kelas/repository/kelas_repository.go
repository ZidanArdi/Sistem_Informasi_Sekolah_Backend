package repository

import (
	"backend/config"
	"backend/modules/kelas/model"

	"gorm.io/gorm"
)

func GetAllKelas(search string, tingkat string) ([]model.Kelas, error) {
	var kelas []model.Kelas

	query := config.DB.Preload("WaliKelas")
	if search != "" {
		query = query.Where("nama_kelas ILIKE ?", "%"+search+"%")
	}
	if tingkat != "" {
		query = query.Where("tingkat = ?", tingkat)
	}

	result := query.Find(&kelas)
	return kelas, result.Error
}

func GetKelasByID(id uint) (model.Kelas, error) {
	var kelas model.Kelas
	result := config.DB.Preload("WaliKelas").First(&kelas, id)
	return kelas, result.Error
}

func CreateKelas(data model.Kelas) (model.Kelas, error) {
	result := config.DB.Create(&data)
	if result.Error == nil {
		config.DB.Preload("WaliKelas").First(&data, data.ID)
	}
	return data, result.Error
}

func UpdateKelas(id uint, data model.Kelas) (model.Kelas, error) {
	var kelas model.Kelas

	if err := config.DB.First(&kelas, id).Error; err != nil {
		return kelas, err
	}

	kelas.NamaKelas = data.NamaKelas
	kelas.Tingkat = data.Tingkat
	kelas.WaliKelasID = data.WaliKelasID

	if err := config.DB.Save(&kelas).Error; err != nil {
		return kelas, err
	}

	config.DB.Preload("WaliKelas").First(&kelas, kelas.ID)
	return kelas, nil
}

func DeleteKelas(id uint) error {
	var kelas model.Kelas

	if err := config.DB.First(&kelas, id).Error; err != nil {
		return err
	}

	return config.DB.Delete(&kelas).Error
}

func IsNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}
