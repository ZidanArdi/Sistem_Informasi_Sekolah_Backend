package repository

import (
	"strconv"

	"backend/config"
	"backend/modules/nilai/model"

	"gorm.io/gorm"
)

func GetAllNilai(siswaID string, mapelID string, semester string, jenisNilai string) ([]model.Nilai, error) {
	var nilai []model.Nilai

	query := config.DB.Preload("Siswa").Preload("Siswa.Kelas").Preload("Mapel")
	if siswaID != "" {
		if parsed, err := strconv.Atoi(siswaID); err == nil {
			query = query.Where("siswa_id = ?", parsed)
		}
	}
	if mapelID != "" {
		if parsed, err := strconv.Atoi(mapelID); err == nil {
			query = query.Where("mapel_id = ?", parsed)
		}
	}
	if semester != "" {
		query = query.Where("semester = ?", semester)
	}
	if jenisNilai != "" {
		query = query.Where("jenis_nilai = ?", jenisNilai)
	}

	result := query.Find(&nilai)
	return nilai, result.Error
}

func GetNilaiByID(id uint) (model.Nilai, error) {
	var nilai model.Nilai
	result := config.DB.Preload("Siswa").Preload("Siswa.Kelas").Preload("Mapel").First(&nilai, id)
	return nilai, result.Error
}

func CreateNilai(data model.Nilai) (model.Nilai, error) {
	result := config.DB.Create(&data)
	if result.Error == nil {
		config.DB.Preload("Siswa").Preload("Siswa.Kelas").Preload("Mapel").First(&data, data.ID)
	}
	return data, result.Error
}

func UpdateNilai(id uint, data model.Nilai) (model.Nilai, error) {
	var nilai model.Nilai

	if err := config.DB.First(&nilai, id).Error; err != nil {
		return nilai, err
	}

	nilai.SiswaID = data.SiswaID
	nilai.MapelID = data.MapelID
	nilai.Semester = data.Semester
	nilai.JenisNilai = data.JenisNilai
	nilai.Nilai = data.Nilai

	if err := config.DB.Save(&nilai).Error; err != nil {
		return nilai, err
	}

	config.DB.Preload("Siswa").Preload("Siswa.Kelas").Preload("Mapel").First(&nilai, nilai.ID)
	return nilai, nil
}

func DeleteNilai(id uint) error {
	var nilai model.Nilai

	if err := config.DB.First(&nilai, id).Error; err != nil {
		return err
	}

	return config.DB.Delete(&nilai).Error
}

func IsNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}
