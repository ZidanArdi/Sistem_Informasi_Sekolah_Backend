package repository

import (
	"strconv"

	"backend/config"
	"backend/modules/nilai/model"

	"gorm.io/gorm"
)

func GetAllNilai(siswaID string, kelasID string, mapelID string, semester string, tahunAjaran string) ([]model.Nilai, error) {
	var nilai []model.Nilai

	query := config.DB.Preload("Siswa").Preload("Siswa.Kelas").Preload("Mapel").Preload("Guru").Preload("Kelas")
	if siswaID != "" {
		if parsed, err := strconv.Atoi(siswaID); err == nil {
			query = query.Where("siswa_id = ?", parsed)
		}
	}
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
	if semester != "" {
		query = query.Where("semester = ?", semester)
	}
	if tahunAjaran != "" {
		query = query.Where("tahun_ajaran = ?", tahunAjaran)
	}

	result := query.Find(&nilai)
	return nilai, result.Error
}

func GetNilaiByID(id uint) (model.Nilai, error) {
	var nilai model.Nilai
	result := config.DB.Preload("Siswa").Preload("Siswa.Kelas").Preload("Mapel").Preload("Guru").Preload("Kelas").First(&nilai, id)
	return nilai, result.Error
}

func CreateNilai(data model.Nilai) (model.Nilai, error) {
	result := config.DB.Create(&data)
	if result.Error == nil {
		config.DB.Preload("Siswa").Preload("Siswa.Kelas").Preload("Mapel").Preload("Guru").Preload("Kelas").First(&data, data.ID)
	}
	return data, result.Error
}

func UpdateNilai(id uint, data model.Nilai) (model.Nilai, error) {
	var nilai model.Nilai

	if err := config.DB.First(&nilai, id).Error; err != nil {
		return nilai, err
	}

	nilai.SiswaID = data.SiswaID
	nilai.KelasID = data.KelasID
	nilai.MapelID = data.MapelID
	nilai.GuruID = data.GuruID
	nilai.Tugas = data.Tugas
	nilai.UTS = data.UTS
	nilai.UAS = data.UAS
	nilai.NilaiAkhir = data.NilaiAkhir
	nilai.GradeHuruf = data.GradeHuruf
	nilai.Semester = data.Semester
	nilai.TahunAjaran = data.TahunAjaran

	if err := config.DB.Save(&nilai).Error; err != nil {
		return nilai, err
	}

	config.DB.Preload("Siswa").Preload("Siswa.Kelas").Preload("Mapel").Preload("Guru").Preload("Kelas").First(&nilai, nilai.ID)
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

func GetSiswaKelasID(siswaID uint) (uint, error) {
	var siswa struct {
		KelasID uint
	}
	err := config.DB.Table("siswas").Select("kelas_id").Where("id = ? AND deleted_at IS NULL", siswaID).Scan(&siswa).Error
	if err != nil || siswa.KelasID == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return siswa.KelasID, nil
}

func CheckExistingNilai(siswaID, mapelID uint, semester, tahunAjaran string) (model.Nilai, error) {
	var existing model.Nilai
	err := config.DB.Where("siswa_id = ? AND mapel_id = ? AND semester = ? AND tahun_ajaran = ?",
		siswaID, mapelID, semester, tahunAjaran).First(&existing).Error
	return existing, err
}

func GetSiswaIDByUserID(userID uint) (uint, error) {
	var siswa struct {
		ID uint
	}
	err := config.DB.Table("siswas").Select("id").Where("user_id = ?", userID).Scan(&siswa).Error
	return siswa.ID, err
}

func SaveNilai(existing model.Nilai, data model.Nilai) (model.Nilai, error) {
	existing.KelasID = data.KelasID
	existing.GuruID = data.GuruID
	existing.Tugas = data.Tugas
	existing.UTS = data.UTS
	existing.UAS = data.UAS
	existing.NilaiAkhir = data.NilaiAkhir
	existing.GradeHuruf = data.GradeHuruf
	
	if errSave := config.DB.Save(&existing).Error; errSave != nil {
		return model.Nilai{}, errSave
	}
	config.DB.Preload("Siswa").Preload("Siswa.Kelas").Preload("Mapel").Preload("Guru").First(&existing, existing.ID)
	return existing, nil
}
