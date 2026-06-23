package repository

import (
	"backend/config"
	"backend/modules/mapel/model"

	"gorm.io/gorm"
)

func GetAllMapel(search string) ([]model.Mapel, error) {
	var mapel []model.Mapel

	query := config.DB
	if search != "" {
		query = query.Where("nama_mapel ILIKE ? OR kode_mapel ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	result := query.Find(&mapel)
	return mapel, result.Error
}

func GetMapelByID(id uint) (model.Mapel, error) {
	var mapel model.Mapel
	result := config.DB.First(&mapel, id)
	return mapel, result.Error
}

func CreateMapel(data model.Mapel) (model.Mapel, error) {
	result := config.DB.Create(&data)
	return data, result.Error
}

func UpdateMapel(id uint, data model.Mapel) (model.Mapel, error) {
	var mapel model.Mapel

	if err := config.DB.First(&mapel, id).Error; err != nil {
		return mapel, err
	}

	mapel.KodeMapel = data.KodeMapel
	mapel.NamaMapel = data.NamaMapel
	mapel.Jam = data.Jam

	err := config.DB.Save(&mapel).Error
	return mapel, err
}

func DeleteMapel(id uint) error {
	var mapel model.Mapel

	if err := config.DB.First(&mapel, id).Error; err != nil {
		return err
	}

	return config.DB.Delete(&mapel).Error
}

func CheckKodeMapelExists(kodeMapel string, excludeID uint) bool {
	var count int64
	query := config.DB.Model(&model.Mapel{}).Where("kode_mapel = ?", kodeMapel)
	if excludeID != 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count > 0
}

func IsNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}
