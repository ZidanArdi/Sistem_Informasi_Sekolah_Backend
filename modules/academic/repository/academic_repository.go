package repository

import (
	"backend/config"
	"backend/modules/guru/model"
	jadwalModel "backend/modules/jadwal/model"
	siswaModel "backend/modules/siswa/model"
)

func GetGuruIDByUserID(userID uint) (uint, error) {
	var guru model.Guru
	if err := config.DB.Select("id").Where("user_id = ? AND deleted_at IS NULL", userID).First(&guru).Error; err != nil {
		return 0, err
	}
	return guru.ID, nil
}

func GetGuruMapelCount(guruID uint, mapelID uint) (int64, error) {
	var count int64
	err := config.DB.Table("guru_mapel").Where("guru_id = ? AND mapel_id = ?", guruID, mapelID).Count(&count).Error
	return count, err
}

func GetJadwalCount(guruID uint, mapelID uint, kelasID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&jadwalModel.Jadwal{}).Where("guru_id = ? AND mapel_id = ? AND kelas_id = ?", guruID, mapelID, kelasID).Count(&count).Error
	return count, err
}

func GetTeacherAssignments(guruID uint) ([]jadwalModel.Jadwal, error) {
	var jadwals []jadwalModel.Jadwal
	err := config.DB.Preload("Kelas").Preload("Mapel").Where("guru_id = ?", guruID).Find(&jadwals).Error
	return jadwals, err
}

func GetTeacherClasses(guruID uint) ([]uint, error) {
	var classIDs []uint
	err := config.DB.Model(&jadwalModel.Jadwal{}).Where("guru_id = ?", guruID).Distinct("kelas_id").Pluck("kelas_id", &classIDs).Error
	return classIDs, err
}

func GetStudentsByTeacher(guruID uint) ([]siswaModel.Siswa, error) {
	var siswas []siswaModel.Siswa
	err := config.DB.Preload("Kelas").Where("kelas_id IN (SELECT DISTINCT kelas_id FROM jadwals WHERE guru_id = ? AND deleted_at IS NULL)", guruID).Find(&siswas).Error
	return siswas, err
}

func GetStudentsByTeacherAndClass(guruID uint, kelasID uint) ([]siswaModel.Siswa, error) {
	var siswas []siswaModel.Siswa
	err := config.DB.Preload("Kelas").Where("kelas_id = ? AND kelas_id IN (SELECT DISTINCT kelas_id FROM jadwals WHERE guru_id = ? AND deleted_at IS NULL)", kelasID, guruID).Find(&siswas).Error
	return siswas, err
}

func CanViewStudent(guruID uint, siswaID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&siswaModel.Siswa{}).Where("id = ? AND kelas_id IN (SELECT DISTINCT kelas_id FROM jadwals WHERE guru_id = ? AND deleted_at IS NULL)", siswaID, guruID).Count(&count).Error
	return count > 0, err
}
