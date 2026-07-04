package repository

import (
	"backend/config"
	"backend/modules/kelas/model"
	siswaModel "backend/modules/siswa/model"

	"gorm.io/gorm"
)

func GetAllKelas(search string, tingkat string) ([]model.Kelas, error) {
	var kelas []model.Kelas

	query := config.DB.Preload("WaliKelas")
	if search != "" {
		query = query.Where("nama_kelas ILIKE ? OR jurusan ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if tingkat != "" {
		query = query.Where("tingkat = ?", tingkat)
	}

	result := query.Find(&kelas)
	if result.Error != nil {
		return nil, result.Error
	}

	for i := range kelas {
		var count int64
		config.DB.Model(&siswaModel.Siswa{}).Where("kelas_id = ?", kelas[i].ID).Count(&count)
		kelas[i].TotalSiswa = int(count)
	}

	return kelas, nil
}

func GetKelasByID(id uint) (model.Kelas, error) {
	var kelas model.Kelas
	result := config.DB.Preload("WaliKelas").First(&kelas, id)
	if result.Error != nil {
		return kelas, result.Error
	}

	var count int64
	config.DB.Model(&siswaModel.Siswa{}).Where("kelas_id = ?", kelas.ID).Count(&count)
	kelas.TotalSiswa = int(count)

	return kelas, nil
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
	kelas.Jurusan = data.Jurusan
	kelas.Kapasitas = data.Kapasitas
	kelas.WaliKelasID = data.WaliKelasID

	if err := config.DB.Save(&kelas).Error; err != nil {
		return kelas, err
	}

	config.DB.Preload("WaliKelas").First(&kelas, kelas.ID)
	
	var count int64
	config.DB.Model(&siswaModel.Siswa{}).Where("kelas_id = ?", kelas.ID).Count(&count)
	kelas.TotalSiswa = int(count)
	
	return kelas, nil
}

func CheckWaliKelasExists(waliKelasID uint, excludeKelasID uint) bool {
	var count int64
	query := config.DB.Unscoped().Model(&model.Kelas{}).Where("wali_kelas_id = ?", waliKelasID)
	if excludeKelasID != 0 {
		query = query.Where("id != ?", excludeKelasID)
	}
	query.Count(&count)
	return count > 0
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
