package repository

import (
	"strconv"

	"backend/config"
	"backend/modules/jadwal/model"

	"gorm.io/gorm"
)

func GetAllJadwal(kelasID string, mapelID string, guruID string, hari string) ([]model.Jadwal, error) {
	var jadwal []model.Jadwal

	query := config.DB.Preload("Kelas").Preload("Kelas.WaliKelas").Preload("Mapel").Preload("Guru")
	if kelasID != "" {
		if parsed, err := strconv.Atoi(kelasID); err == nil {
			query = query.Where("kelas_id = ?", parsed)
		}
	}
	if mapelID != "" {
		if parsed, err := strconv.Atoi(mapelID); err == nil {
			query = query.Where("mapel_id = ?", parsed)
		}
	}
	if guruID != "" {
		if parsed, err := strconv.Atoi(guruID); err == nil {
			query = query.Where("guru_id = ?", parsed)
		}
	}
	if hari != "" {
		query = query.Where("hari ILIKE ?", hari)
	}

	result := query.Find(&jadwal)
	return jadwal, result.Error
}

func GetJadwalByID(id uint) (model.Jadwal, error) {
	var jadwal model.Jadwal
	result := config.DB.Preload("Kelas").Preload("Kelas.WaliKelas").Preload("Mapel").Preload("Guru").First(&jadwal, id)
	return jadwal, result.Error
}

func CreateJadwal(data model.Jadwal) (model.Jadwal, error) {
	result := config.DB.Create(&data)
	if result.Error == nil {
		config.DB.Preload("Kelas").Preload("Kelas.WaliKelas").Preload("Mapel").Preload("Guru").First(&data, data.ID)
	}
	return data, result.Error
}

func UpdateJadwal(id uint, data model.Jadwal) (model.Jadwal, error) {
	var jadwal model.Jadwal

	if err := config.DB.First(&jadwal, id).Error; err != nil {
		return jadwal, err
	}

	jadwal.KelasID = data.KelasID
	jadwal.MapelID = data.MapelID
	jadwal.GuruID = data.GuruID
	jadwal.Hari = data.Hari
	jadwal.JamMulai = data.JamMulai
	jadwal.JamSelesai = data.JamSelesai

	if err := config.DB.Save(&jadwal).Error; err != nil {
		return jadwal, err
	}

	config.DB.Preload("Kelas").Preload("Kelas.WaliKelas").Preload("Mapel").Preload("Guru").First(&jadwal, jadwal.ID)
	return jadwal, nil
}

func DeleteJadwal(id uint) error {
	var jadwal model.Jadwal

	if err := config.DB.First(&jadwal, id).Error; err != nil {
		return err
	}

	return config.DB.Delete(&jadwal).Error
}

func IsNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}
