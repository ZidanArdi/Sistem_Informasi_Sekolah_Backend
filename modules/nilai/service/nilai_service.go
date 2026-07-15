package service

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"backend/config"
	academicModel "backend/modules/academic/model"
	academicService "backend/modules/academic/service"
	"backend/modules/nilai/model"
	"backend/modules/nilai/repository"
)

func GetAllNilai(ctx academicModel.TeacherContext, querySiswaID string, queryKelasID string, queryMapelID string, querySemester string, queryTahunAjaran string) ([]model.Nilai, error) {
	if ctx.Role == "siswa" {
		siswaID, err := repository.GetSiswaIDByUserID(ctx.UserID)
		if err != nil || siswaID == 0 {
			return nil, errors.New("data siswa tidak ditemukan")
		}
		querySiswaID = strconv.Itoa(int(siswaID))
	}

	return repository.GetAllNilai(querySiswaID, queryKelasID, queryMapelID, querySemester, queryTahunAjaran)
}

func GetNilaiByID(id uint) (model.Nilai, error) {
	return repository.GetNilaiByID(id)
}

func CreateNilai(ctx academicModel.TeacherContext, data model.Nilai) (model.Nilai, error) {
	// 1. Resolve student's KelasID
	kelasID, err := repository.GetSiswaKelasID(data.SiswaID)
	if err != nil {
		return model.Nilai{}, errors.New("siswa tidak ditemukan")
	}
	data.KelasID = kelasID

	// 2. Validate authorization if user is a Guru
	if ctx.Role == "guru" {
		data.GuruID = ctx.GuruID
		errAuth := academicService.CanInputGrade(ctx, data.MapelID, data.KelasID)
		if errAuth != nil {
			return model.Nilai{}, errAuth
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

	// Debug Logging
	if config.Debug {
		log.Printf("[DEBUG] CreateNilai - UserID: %d, GuruID: %d, MapelID: %d, KelasID: %d, SiswaID: %d, Sem: %s, TA: %s",
			ctx.UserID, ctx.GuruID, data.MapelID, data.KelasID, data.SiswaID, data.Semester, data.TahunAjaran)
	}

	// 5. Check if record already exists (UPSERT check)
	existing, errCheck := repository.CheckExistingNilai(data.SiswaID, data.MapelID, data.Semester, data.TahunAjaran)
	if errCheck == nil {
		// Update existing
		return repository.SaveNilai(existing, data)
	}

	// Insert new
	return repository.CreateNilai(data)
}

func UpdateNilai(id uint, ctx academicModel.TeacherContext, data model.Nilai) (model.Nilai, error) {
	// Resolve student's KelasID
	kelasID, err := repository.GetSiswaKelasID(data.SiswaID)
	if err != nil {
		return model.Nilai{}, errors.New("siswa tidak ditemukan")
	}
	data.KelasID = kelasID

	// Validate authorization
	if ctx.Role == "guru" {
		data.GuruID = ctx.GuruID
		errAuth := academicService.CanInputGrade(ctx, data.MapelID, data.KelasID)
		if errAuth != nil {
			return model.Nilai{}, errAuth
		}
	} else {
		return model.Nilai{}, errors.New("Hanya guru yang dapat menginput atau mengedit nilai.")
	}

	data.NilaiAkhir = (data.Tugas * 0.3) + (data.UTS * 0.3) + (data.UAS * 0.4)
	data.GradeHuruf = deriveGradeHuruf(data.NilaiAkhir)

	if data.Tugas < 0 || data.Tugas > 100 || data.UTS < 0 || data.UTS > 100 || data.UAS < 0 || data.UAS > 100 {
		return model.Nilai{}, errors.New("nilai tugas, UTS, dan UAS harus di antara 0 dan 100")
	}

	if config.Debug {
		log.Printf("[DEBUG] UpdateNilai - NilaiID: %d, UserID: %d, GuruID: %d, MapelID: %d, KelasID: %d, SiswaID: %d",
			id, ctx.UserID, ctx.GuruID, data.MapelID, data.KelasID, data.SiswaID)
	}

	return repository.UpdateNilai(id, data)
}

func DeleteNilai(id uint, ctx academicModel.TeacherContext) error {
	if ctx.Role != "guru" {
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
