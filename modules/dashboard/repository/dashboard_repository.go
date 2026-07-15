package repository

import (
	"backend/config"
	authModel "backend/modules/auth/model"
	guruModel "backend/modules/guru/model"
	jadwalModel "backend/modules/jadwal/model"
	kelasModel "backend/modules/kelas/model"
	mapelModel "backend/modules/mapel/model"
	nilaiModel "backend/modules/nilai/model"
	perizinanModel "backend/modules/perizinan/model"
	siswaModel "backend/modules/siswa/model"
)

func GetUserByID(userID uint) (authModel.User, error) {
	var user authModel.User
	err := config.DB.First(&user, userID).Error
	return user, err
}

func CountSiswa() (int64, error) {
	var count int64
	err := config.DB.Model(&siswaModel.Siswa{}).Count(&count).Error
	return count, err
}

func CountGuru() (int64, error) {
	var count int64
	err := config.DB.Model(&guruModel.Guru{}).Count(&count).Error
	return count, err
}

func CountKelas() (int64, error) {
	var count int64
	err := config.DB.Model(&kelasModel.Kelas{}).Count(&count).Error
	return count, err
}

func CountMapel() (int64, error) {
	var count int64
	err := config.DB.Model(&mapelModel.Mapel{}).Count(&count).Error
	return count, err
}

func CountPerizinanPending() (int64, error) {
	var count int64
	err := config.DB.Model(&perizinanModel.Perizinan{}).Where("status = ?", "Pending").Count(&count).Error
	return count, err
}

func CountTodaySchedule(guruID uint, dayName string) (int64, error) {
	var count int64
	err := config.DB.Model(&jadwalModel.Jadwal{}).Where("guru_id = ? AND hari = ?", guruID, dayName).Count(&count).Error
	return count, err
}

func CountTotalClasses(guruID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&jadwalModel.Jadwal{}).Where("guru_id = ?", guruID).Distinct("kelas_id").Count(&count).Error
	return count, err
}

func CountTotalStudentsForGuru(guruID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&siswaModel.Siswa{}).Where("kelas_id IN (SELECT DISTINCT kelas_id FROM jadwals WHERE guru_id = ? AND deleted_at IS NULL)", guruID).Count(&count).Error
	return count, err
}

func GetJadwalByGuru(guruID uint) ([]jadwalModel.Jadwal, error) {
	var schedules []jadwalModel.Jadwal
	err := config.DB.Where("guru_id = ?", guruID).Find(&schedules).Error
	return schedules, err
}

func CountSiswaByKelas(kelasID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&siswaModel.Siswa{}).Where("kelas_id = ?", kelasID).Count(&count).Error
	return count, err
}

func CountNilaiByJadwal(mapelID uint, kelasID uint, semester string, tahunAjaran string) (int64, error) {
	var count int64
	err := config.DB.Model(&nilaiModel.Nilai{}).Where("mapel_id = ? AND kelas_id = ? AND semester = ? AND tahun_ajaran = ?", mapelID, kelasID, semester, tahunAjaran).Count(&count).Error
	return count, err
}

func GetTodayJadwal(guruID uint, dayName string) ([]jadwalModel.Jadwal, error) {
	var jadwals []jadwalModel.Jadwal
	err := config.DB.Preload("Kelas").Preload("Mapel").Where("guru_id = ? AND hari = ?", guruID, dayName).Find(&jadwals).Error
	return jadwals, err
}

func GetSiswaByUserID(userID uint) (siswaModel.Siswa, error) {
	var siswa siswaModel.Siswa
	err := config.DB.Preload("Kelas").Where("user_id = ?", userID).First(&siswa).Error
	return siswa, err
}

func GetTodayJadwalForKelas(kelasID uint, dayName string) ([]jadwalModel.Jadwal, error) {
	var jadwals []jadwalModel.Jadwal
	err := config.DB.Preload("Guru").Preload("Mapel").Where("kelas_id = ? AND hari = ?", kelasID, dayName).Find(&jadwals).Error
	return jadwals, err
}

func GetRiwayatPerizinanSiswa(siswaID uint, limit int) ([]perizinanModel.Perizinan, error) {
	var perizinans []perizinanModel.Perizinan
	err := config.DB.Preload("Guru").Where("siswa_id = ?", siswaID).Order("created_at desc").Limit(limit).Find(&perizinans).Error
	return perizinans, err
}
