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
	perizinanModel "backend/modules/perizinan/model"
	siswaModel "backend/modules/siswa/model"
)

type AdminDashboardData struct {
	TotalSiswa            int64 `json:"total_siswa"`
	TotalGuru             int64 `json:"total_guru"`
	TotalKelas            int64 `json:"total_kelas"`
	TotalMapel            int64 `json:"total_mapel"`
	TotalPerizinanPending int64 `json:"total_perizinan_pending"`
}

type GuruDashboardData struct {
	TotalSiswa            int64                 `json:"total_siswa"`
	TotalPerizinanPending int64                 `json:"total_perizinan_pending"`
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

func GetAdminDashboard(role string) (AdminDashboardData, error) {
	role = strings.ToLower(strings.TrimSpace(role))
	if role != "admin" {
		return AdminDashboardData{}, errors.New("akses ditolak: hanya admin yang dapat mengakses dasbor ini")
	}

	var data AdminDashboardData

	config.DB.Model(&siswaModel.Siswa{}).Count(&data.TotalSiswa)
	config.DB.Model(&guruModel.Guru{}).Count(&data.TotalGuru)
	config.DB.Model(&kelasModel.Kelas{}).Count(&data.TotalKelas)
	config.DB.Model(&mapelModel.Mapel{}).Count(&data.TotalMapel)
	config.DB.Model(&perizinanModel.Perizinan{}).Where("status = ?", "Pending").Count(&data.TotalPerizinanPending)

	return data, nil
}

func GetGuruDashboard(role string, email string) (GuruDashboardData, error) {
	role = strings.ToLower(strings.TrimSpace(role))
	if role != "guru" {
		return GuruDashboardData{}, errors.New("akses ditolak: hanya guru yang dapat mengakses dasbor ini")
	}

	var data GuruDashboardData
	todayDayName := getIndonesianDayName()

	var user authModel.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return GuruDashboardData{}, errors.New("akun guru tidak ditemukan")
	}

	var guru guruModel.Guru
	if err := config.DB.Where("user_id = ?", user.ID).First(&guru).Error; err != nil {
		return GuruDashboardData{}, errors.New("data profil guru tidak ditemukan")
	}

	// 1. Total Siswa (School wide)
	config.DB.Model(&siswaModel.Siswa{}).Count(&data.TotalSiswa)

	// 2. Total Perizinan Pending (Siswa -> Semua Guru approval)
	config.DB.Model(&perizinanModel.Perizinan{}).Where("status = ?", "Pending").Count(&data.TotalPerizinanPending)

	// 3. Jadwal Hari Ini
	config.DB.Preload("Kelas").Preload("Mapel").Where("guru_id = ? AND hari = ?", guru.ID, todayDayName).Find(&data.JadwalHariIni)

	return data, nil
}

func GetSiswaDashboard(role string, email string) (SiswaDashboardData, error) {
	role = strings.ToLower(strings.TrimSpace(role))
	if role != "siswa" {
		return SiswaDashboardData{}, errors.New("akses ditolak: hanya siswa yang dapat mengakses dasbor ini")
	}

	var data SiswaDashboardData
	todayDayName := getIndonesianDayName()

	// Get Siswa Details with Class
	if err := config.DB.Preload("Kelas").Where("nis = ?", email).First(&data.ProfilSiswa).Error; err != nil {
		return SiswaDashboardData{}, errors.New("data siswa tidak ditemukan")
	}

	// 1. Jadwal Hari Ini untuk kelas siswa tersebut
	config.DB.Preload("Guru").Preload("Mapel").Where("kelas_id = ? AND hari = ?", data.ProfilSiswa.KelasID, todayDayName).Find(&data.JadwalHariIni)

	// 2. Riwayat Perizinan Siswa (Limit 5)
	config.DB.Preload("Guru").Where("siswa_id = ?", data.ProfilSiswa.ID).Order("created_at desc").Limit(5).Find(&data.RiwayatPerizinan)

	return data, nil
}
