package service

import (
	"errors"
	"strconv"
	"strings"

	"backend/config"
	"backend/modules/nilai/model"
	"backend/modules/nilai/repository"
)

func GetAllNilai(userID uint, role string, querySiswaID string, queryKelasID string, queryMapelID string, querySemester string, queryTahunAjaran string) ([]model.Nilai, error) {
	// Enforce Siswa role restriction to only see their own grades
	if role == "siswa" {
		var siswa struct {
			ID uint
		}
		if err := config.DB.Table("siswas").Select("id").Where("user_id = ?", userID).Scan(&siswa).Error; err != nil || siswa.ID == 0 {
			return nil, errors.New("data siswa tidak ditemukan")
		}
		querySiswaID = strconv.Itoa(int(siswa.ID))
	}

	return repository.GetAllNilai(querySiswaID, queryKelasID, queryMapelID, querySemester, queryTahunAjaran)
}

func GetNilaiByID(id uint) (model.Nilai, error) {
	return repository.GetNilaiByID(id)
}

func CreateNilai(userID uint, role string, data model.Nilai) (model.Nilai, error) {
	// 1. Resolve student's KelasID
	var siswa struct {
		KelasID uint
	}
	if err := config.DB.Table("siswas").Select("kelas_id").Where("id = ? AND deleted_at IS NULL", data.SiswaID).Scan(&siswa).Error; err != nil || data.SiswaID == 0 {
		return model.Nilai{}, errors.New("siswa tidak ditemukan")
	}
	data.KelasID = siswa.KelasID

	// 2. Resolve GuruID and validate authorization if user is a Guru
	if role == "guru" {
		var guru struct {
			ID uint
		}
		if err := config.DB.Table("gurus").Select("id").Where("user_id = ? AND deleted_at IS NULL", userID).Scan(&guru).Error; err != nil || guru.ID == 0 {
			return model.Nilai{}, errors.New("data guru tidak ditemukan")
		}
		data.GuruID = guru.ID

		// Check eligibility in guru_mapel
		var mapelCount int64
		config.DB.Table("guru_mapels").Where("guru_id = ? AND mapel_id = ?", guru.ID, data.MapelID).Count(&mapelCount)
		if mapelCount == 0 {
			return model.Nilai{}, errors.New("Anda tidak memiliki akses untuk menginput nilai pada kelas ini.")
		}

		// Check schedule in jadwal
		var jadwalCount int64
		config.DB.Table("jadwals").Where("guru_id = ? AND mapel_id = ? AND kelas_id = ? AND deleted_at IS NULL", guru.ID, data.MapelID, data.KelasID).Count(&jadwalCount)
		if jadwalCount == 0 {
			return model.Nilai{}, errors.New("Anda tidak memiliki akses untuk menginput nilai pada kelas ini.")
		}
	} else {
		return model.Nilai{}, errors.New("Hanya guru yang dapat menginput atau mengedit nilai.")
	}

	// 3. Compute final score and letter grade
	data.NilaiAkhir = (data.Tugas * 0.3) + (data.UTS * 0.3) + (data.UAS * 0.4)
	data.GradeHuruf = deriveGradeHuruf(data.NilaiAkhir)

	// 4. Validate input values
	if data.Tugas < 0 || data.Tugas > 100 || data.UTS < 0 || data.UTS > 100 || data.UAS < 0 || data.UAS > 100 {
		return model.Nilai{}, errors.New("nilai tugas, UTS, dan UAS harus di antara 0 dan 100")
	}
	if strings.TrimSpace(data.Semester) == "" || strings.TrimSpace(data.TahunAjaran) == "" {
		return model.Nilai{}, errors.New("semester dan tahun ajaran wajib diisi")
	}

	// 5. Check if record already exists (UPSERT check)
	var existing model.Nilai
	err := config.DB.Where("siswa_id = ? AND mapel_id = ? AND semester = ? AND tahun_ajaran = ?",
		data.SiswaID, data.MapelID, data.Semester, data.TahunAjaran).First(&existing).Error
	
	if err == nil {
		// Update existing
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

	// Insert new
	created, errCreate := repository.CreateNilai(data)
	return created, errCreate
}

func UpdateNilai(id uint, userID uint, role string, data model.Nilai) (model.Nilai, error) {
	// Resolve student's KelasID
	var siswa struct {
		KelasID uint
	}
	if err := config.DB.Table("siswas").Select("kelas_id").Where("id = ? AND deleted_at IS NULL", data.SiswaID).Scan(&siswa).Error; err != nil || data.SiswaID == 0 {
		return model.Nilai{}, errors.New("siswa tidak ditemukan")
	}
	data.KelasID = siswa.KelasID

	// Resolve Guru ID and validate authorization
	if role == "guru" {
		var guru struct {
			ID uint
		}
		if err := config.DB.Table("gurus").Select("id").Where("user_id = ? AND deleted_at IS NULL", userID).Scan(&guru).Error; err != nil || guru.ID == 0 {
			return model.Nilai{}, errors.New("data guru tidak ditemukan")
		}
		data.GuruID = guru.ID

		// Check eligibility
		var mapelCount int64
		config.DB.Table("guru_mapels").Where("guru_id = ? AND mapel_id = ?", guru.ID, data.MapelID).Count(&mapelCount)
		if mapelCount == 0 {
			return model.Nilai{}, errors.New("Anda tidak memiliki akses untuk menginput nilai pada kelas ini.")
		}

		// Check schedule
		var jadwalCount int64
		config.DB.Table("jadwals").Where("guru_id = ? AND mapel_id = ? AND kelas_id = ? AND deleted_at IS NULL", guru.ID, data.MapelID, data.KelasID).Count(&jadwalCount)
		if jadwalCount == 0 {
			return model.Nilai{}, errors.New("Anda tidak memiliki akses untuk menginput nilai pada kelas ini.")
		}
	} else {
		return model.Nilai{}, errors.New("Hanya guru yang dapat menginput atau mengedit nilai.")
	}

	data.NilaiAkhir = (data.Tugas * 0.3) + (data.UTS * 0.3) + (data.UAS * 0.4)
	data.GradeHuruf = deriveGradeHuruf(data.NilaiAkhir)

	if data.Tugas < 0 || data.Tugas > 100 || data.UTS < 0 || data.UTS > 100 || data.UAS < 0 || data.UAS > 100 {
		return model.Nilai{}, errors.New("nilai tugas, UTS, dan UAS harus di antara 0 dan 100")
	}

	return repository.UpdateNilai(id, data)
}

func DeleteNilai(id uint, role string) error {
	if role != "guru" {
		return errors.New("Hanya guru yang dapat menghapus nilai.")
	}
	return repository.DeleteNilai(id)
}

func deriveGradeHuruf(score float64) string {
	if score >= 85 {
		return "A"
	}
	if score >= 75 {
		return "B"
	}
	if score >= 65 {
		return "C"
	}
	return "D"
}
