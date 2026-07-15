package service

import (
	"errors"

	"backend/modules/academic/model"
	"backend/modules/academic/repository"
	jadwalModel "backend/modules/jadwal/model"
	siswaModel "backend/modules/siswa/model"
)

// GetTeacherAssignments returns all schedules for a given teacher context
func GetTeacherAssignments(ctx model.TeacherContext) ([]jadwalModel.Jadwal, error) {
	if ctx.Role != "guru" {
		return nil, errors.New("ERR_FORBIDDEN: hanya guru yang memiliki assignment mengajar")
	}
	return repository.GetTeacherAssignments(ctx.GuruID)
}

// GetTeacherClasses returns a list of unique class IDs taught by the teacher
func GetTeacherClasses(ctx model.TeacherContext) ([]uint, error) {
	if ctx.Role != "guru" {
		return nil, errors.New("ERR_FORBIDDEN: hanya guru yang memiliki assignment kelas")
	}
	return repository.GetTeacherClasses(ctx.GuruID)
}

// GetStudentsByTeacher returns all students in all classes taught by the teacher
func GetStudentsByTeacher(ctx model.TeacherContext) ([]siswaModel.Siswa, error) {
	if ctx.Role != "guru" {
		return nil, errors.New("ERR_FORBIDDEN: akses ditolak")
	}
	return repository.GetStudentsByTeacher(ctx.GuruID)
}

// GetStudentsByTeacherAndClass returns all students in a specific class taught by the teacher
func GetStudentsByTeacherAndClass(ctx model.TeacherContext, kelasID uint) ([]siswaModel.Siswa, error) {
	if ctx.Role != "guru" {
		return nil, errors.New("ERR_FORBIDDEN: akses ditolak")
	}
	return repository.GetStudentsByTeacherAndClass(ctx.GuruID, kelasID)
}

// CanInputGrade checks if a teacher is authorized to input grades for a specific subject and class
func CanInputGrade(ctx model.TeacherContext, mapelID uint, kelasID uint) error {
	if ctx.Role != "guru" {
		return errors.New("Hanya guru yang dapat menginput atau mengedit nilai.")
	}

	// 1. Check eligibility in guru_mapel
	mapelCount, err := repository.GetGuruMapelCount(ctx.GuruID, mapelID)
	if err != nil {
		return errors.New("ERR_INTERNAL_SERVER: gagal memeriksa akses mapel")
	}
	if mapelCount == 0 {
		return errors.New("Anda tidak memiliki akses untuk menginput nilai pada kelas ini.")
	}

	// 2. Check schedule in jadwals
	jadwalCount, err := repository.GetJadwalCount(ctx.GuruID, mapelID, kelasID)
	if err != nil {
		return errors.New("ERR_INTERNAL_SERVER: gagal memeriksa akses jadwal")
	}
	if jadwalCount == 0 {
		return errors.New("Anda tidak memiliki akses untuk menginput nilai pada kelas ini.")
	}

	return nil
}

// CanViewStudent checks if a teacher is authorized to view a specific student's data
func CanViewStudent(ctx model.TeacherContext, siswaID uint) error {
	if ctx.Role != "guru" {
		return errors.New("ERR_FORBIDDEN: akses ditolak")
	}
	authorized, err := repository.CanViewStudent(ctx.GuruID, siswaID)
	if err != nil {
		return errors.New("ERR_INTERNAL_SERVER: gagal memeriksa otoritas data siswa")
	}
	if !authorized {
		return errors.New("ERR_FORBIDDEN: Anda tidak memiliki akses ke data siswa ini")
	}
	return nil
}
