package service

import (
	"errors"
	"strings"
	"time"

	"backend/config"
	authModel "backend/modules/auth/model"
	guruModel "backend/modules/guru/model"
	jadwalModel "backend/modules/jadwal/model"
	kelasModel "backend/modules/kelas/model"
	mapelModel "backend/modules/mapel/model"
	nilaiModel "backend/modules/nilai/model"
	perizinanModel "backend/modules/perizinan/model"
	siswaModel "backend/modules/siswa/model"

	"gorm.io/gorm"
)

type AdminDashboardData struct {
	TotalSiswa            int64 `json:"total_siswa"`
	TotalGuru             int64 `json:"total_guru"`
	TotalKelas            int64 `json:"total_kelas"`
	TotalMapel            int64 `json:"total_mapel"`
	TotalPerizinanPending int64 `json:"total_perizinan_pending"`
}

type GuruDashboardData struct {
	TodaySchedule         int64                 `json:"today_schedule"`
	TotalClasses          int64                 `json:"total_classes"`
	TotalStudents         int64                 `json:"total_students"`
	PendingGrades         int64                 `json:"pending_grades"`
	JadwalHariIni         []jadwalModel.Jadwal  `json:"jadwal_hari_ini"`
}

type SiswaDashboardData struct {
	ProfilSiswa      siswaModel.Siswa            `json:"profil_siswa"`
	JadwalHariIni    []jadwalModel.Jadwal        `json:"jadwal_hari_ini"`
	RiwayatPerizinan []perizinanModel.Perizinan  `json:"riwayat_perizinan"`
}

func getIndonesianDayName() string {
	day := time.Now().Weekday()
	switch day {
	case time.Sunday:
		return "Minggu"
	case time.Monday:
		return "Senin"
	case time.Tuesday:
		return "Selasa"
	case time.Wednesday:
		return "Rabu"
	case time.Thursday:
		return "Kamis"
	case time.Friday:
		return "Jumat"
	case time.Saturday:
		return "Sabtu"
	default:
		return ""
	}
}

func GetAdminDashboard(userID uint) (AdminDashboardData, error) {
	var user authModel.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return AdminDashboardData{}, errors.New("ERR_PROFILE_NOT_FOUND: akun tidak ditemukan")
		}
		return AdminDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan database")
	}

	if strings.ToLower(user.Role) != "admin" {
		return AdminDashboardData{}, errors.New("ERR_FORBIDDEN: akses hanya untuk admin")
	}

	var data AdminDashboardData

	if err := config.DB.Model(&siswaModel.Siswa{}).Count(&data.TotalSiswa).Error; err != nil {
		return AdminDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat data total siswa")
	}
	if err := config.DB.Model(&guruModel.Guru{}).Count(&data.TotalGuru).Error; err != nil {
		return AdminDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat data total guru")
	}
	if err := config.DB.Model(&kelasModel.Kelas{}).Count(&data.TotalKelas).Error; err != nil {
		return AdminDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat data total kelas")
	}
	if err := config.DB.Model(&mapelModel.Mapel{}).Count(&data.TotalMapel).Error; err != nil {
		return AdminDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat data total mata pelajaran")
	}
	if err := config.DB.Model(&perizinanModel.Perizinan{}).Where("status = ?", "Pending").Count(&data.TotalPerizinanPending).Error; err != nil {
		return AdminDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat data total perizinan pending")
	}

	return data, nil
}

func GetGuruDashboard(userID uint) (GuruDashboardData, error) {
	var user authModel.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return GuruDashboardData{}, errors.New("ERR_PROFILE_NOT_FOUND: akun tidak ditemukan")
		}
		return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan database")
	}

	if strings.ToLower(user.Role) != "guru" {
		return GuruDashboardData{}, errors.New("ERR_FORBIDDEN: akses hanya untuk guru")
	}

	var guru guruModel.Guru
	if err := config.DB.Where("user_id = ?", user.ID).First(&guru).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return GuruDashboardData{}, errors.New("ERR_PROFILE_NOT_FOUND: data profil guru tidak ditemukan")
		}
		return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan database")
	}

	var data GuruDashboardData
	todayDayName := getIndonesianDayName()

	// 1. Today's Schedule Count
	if err := config.DB.Model(&jadwalModel.Jadwal{}).Where("guru_id = ? AND hari = ?", guru.ID, todayDayName).Count(&data.TodaySchedule).Error; err != nil {
		return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat jadwal")
	}

	// 2. Total Classes Assigned
	if err := config.DB.Model(&jadwalModel.Jadwal{}).Where("guru_id = ?", guru.ID).Distinct("kelas_id").Count(&data.TotalClasses).Error; err != nil {
		return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat kelas")
	}

	// 3. Total Students in Classes Taught
	if err := config.DB.Model(&siswaModel.Siswa{}).Where("kelas_id IN (SELECT DISTINCT kelas_id FROM jadwals WHERE guru_id = ? AND deleted_at IS NULL)", guru.ID).Count(&data.TotalStudents).Error; err != nil {
		return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat total siswa")
	}

	// 4. Pending Grade Completion
	var schedules []jadwalModel.Jadwal
	if err := config.DB.Where("guru_id = ?", guru.ID).Find(&schedules).Error; err != nil {
		return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat jadwal perwalian")
	}

	var pendingGrades int64 = 0
	for _, sched := range schedules {
		var classStudentsCount int64
		if err := config.DB.Model(&siswaModel.Siswa{}).Where("kelas_id = ?", sched.KelasID).Count(&classStudentsCount).Error; err != nil {
			return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat data siswa kelas")
		}

		var gradedCount int64
		if err := config.DB.Model(&nilaiModel.Nilai{}).
			Where("mapel_id = ? AND kelas_id = ? AND semester = ? AND tahun_ajaran = ?", sched.MapelID, sched.KelasID, sched.Semester, sched.TahunAjaran).
			Count(&gradedCount).Error; err != nil {
			return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat nilai siswa")
		}

		if classStudentsCount > gradedCount {
			pendingGrades += (classStudentsCount - gradedCount)
		}
	}
	data.PendingGrades = pendingGrades

	// 5. Today's schedules list
	if err := config.DB.Preload("Kelas").Preload("Mapel").Where("guru_id = ? AND hari = ?", guru.ID, todayDayName).Find(&data.JadwalHariIni).Error; err != nil {
		return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat daftar jadwal hari ini")
	}

	// Populate TotalSiswa for each schedule's class
	for i := range data.JadwalHariIni {
		var count int64
		if err := config.DB.Model(&siswaModel.Siswa{}).Where("kelas_id = ?", data.JadwalHariIni[i].KelasID).Count(&count).Error; err != nil {
			return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat total siswa per kelas")
		}
		data.JadwalHariIni[i].Kelas.TotalSiswa = int(count)
	}

	return data, nil
}

func GetSiswaDashboard(userID uint) (SiswaDashboardData, error) {
	var user authModel.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return SiswaDashboardData{}, errors.New("ERR_PROFILE_NOT_FOUND: akun tidak ditemukan")
		}
		return SiswaDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan database")
	}

	if strings.ToLower(user.Role) != "siswa" {
		return SiswaDashboardData{}, errors.New("ERR_FORBIDDEN: akses hanya untuk siswa")
	}

	var data SiswaDashboardData
	todayDayName := getIndonesianDayName()

	// Get Siswa Details with Class
	if err := config.DB.Preload("Kelas").Where("user_id = ?", user.ID).First(&data.ProfilSiswa).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return SiswaDashboardData{}, errors.New("ERR_PROFILE_NOT_FOUND: data profil siswa tidak ditemukan")
		}
		return SiswaDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan database")
	}

	// 1. Jadwal Hari Ini untuk kelas siswa tersebut
	if err := config.DB.Preload("Guru").Preload("Mapel").Where("kelas_id = ? AND hari = ?", data.ProfilSiswa.KelasID, todayDayName).Find(&data.JadwalHariIni).Error; err != nil {
		return SiswaDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat jadwal pelajaran hari ini")
	}

	// 2. Riwayat Perizinan Siswa (Limit 5)
	if err := config.DB.Preload("Guru").Where("siswa_id = ?", data.ProfilSiswa.ID).Order("created_at desc").Limit(5).Find(&data.RiwayatPerizinan).Error; err != nil {
		return SiswaDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat riwayat perizinan")
	}

	return data, nil
}

