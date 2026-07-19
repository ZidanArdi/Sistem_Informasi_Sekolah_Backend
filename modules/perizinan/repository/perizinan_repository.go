package repository

import (
	"backend/config"
	"backend/modules/perizinan/model"
	guruModel "backend/modules/guru/model"
	siswaModel "backend/modules/siswa/model"

	"gorm.io/gorm"
)

func CreatePerizinan(data model.Perizinan) (model.Perizinan, error) {
	result := config.DB.Create(&data)
	if result.Error == nil {
		// Preload relations for response consistency
		config.DB.Preload("Siswa").Preload("Siswa.Kelas").Preload("Guru").First(&data, data.ID)
	}
	return data, result.Error
}

func GetPerizinanByID(id uint) (model.Perizinan, error) {
	var perizinan model.Perizinan
	result := config.DB.Preload("Siswa").Preload("Siswa.Kelas").Preload("Guru").First(&perizinan, id)
	return perizinan, result.Error
}

func GetAllPerizinan() ([]model.Perizinan, error) {
	var list []model.Perizinan
	result := config.DB.Preload("Siswa").Preload("Siswa.Kelas").Preload("Guru").
		Joins("JOIN siswas ON siswas.id = perizinan.siswa_id AND siswas.deleted_at IS NULL").
		Order("perizinan.created_at desc").Find(&list)
	return list, result.Error
}

func GetPerizinanBySiswaID(siswaID uint) ([]model.Perizinan, error) {
	var list []model.Perizinan
	result := config.DB.Preload("Siswa").Preload("Siswa.Kelas").Preload("Guru").
		Joins("JOIN siswas ON siswas.id = perizinan.siswa_id AND siswas.deleted_at IS NULL").
		Where("perizinan.siswa_id = ?", siswaID).Order("perizinan.created_at desc").Find(&list)
	return list, result.Error
}

func GetPendingPerizinan() ([]model.Perizinan, error) {
	var list []model.Perizinan
	result := config.DB.Preload("Siswa").Preload("Siswa.Kelas").Preload("Guru").
		Joins("JOIN siswas ON siswas.id = perizinan.siswa_id AND siswas.deleted_at IS NULL").
		Where("perizinan.status = ?", "Pending").Order("perizinan.created_at desc").Find(&list)
	return list, result.Error
}

func UpdatePerizinan(id uint, data model.Perizinan) (model.Perizinan, error) {
	var perizinan model.Perizinan
	if err := config.DB.First(&perizinan, id).Error; err != nil {
		return perizinan, err
	}

	perizinan.Status = data.Status
	perizinan.KeteranganGuru = data.KeteranganGuru
	perizinan.DisetujuiOleh = data.DisetujuiOleh

	err := config.DB.Save(&perizinan).Error
	if err == nil {
		config.DB.Preload("Siswa").Preload("Siswa.Kelas").Preload("Guru").First(&perizinan, perizinan.ID)
	}
	return perizinan, err
}

func GetSiswaIDByEmail(email string) (uint, error) {
	var siswa siswaModel.Siswa
	err := config.DB.Select("id").Where("nis = ?", email).First(&siswa).Error
	if err != nil {
		return 0, err
	}
	return siswa.ID, nil
}

func GetGuruIDByEmail(email string) (uint, error) {
	var user struct {
		ID uint
	}
	err := config.DB.Table("users").Select("id").Where("email = ?", email).First(&user).Error
	if err != nil {
		return 0, err
	}

	var guru guruModel.Guru
	err = config.DB.Select("id").Where("user_id = ?", user.ID).First(&guru).Error
	if err != nil {
		return 0, err
	}
	return guru.ID, nil
}

func IsNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}
