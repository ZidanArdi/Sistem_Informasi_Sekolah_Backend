package service

import (
	"errors"
	"strings"
	"time"

	academicModel "backend/modules/academic/model"
	"backend/modules/dashboard/repository"
	jadwalModel "backend/modules/jadwal/model"
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
	user, err := repository.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return AdminDashboardData{}, errors.New("ERR_PROFILE_NOT_FOUND: akun tidak ditemukan")
		}
		return AdminDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan database")
	}

	if strings.ToLower(user.Role) != "admin" {
		return AdminDashboardData{}, errors.New("ERR_FORBIDDEN: akses hanya untuk admin")
	}

	var data AdminDashboardData

	data.TotalSiswa, err = repository.CountSiswa()
	if err != nil {
		return AdminDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat data total siswa")
	}
	data.TotalGuru, err = repository.CountGuru()
	if err != nil {
		return AdminDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat data total guru")
	}
	data.TotalKelas, err = repository.CountKelas()
	if err != nil {
		return AdminDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat data total kelas")
	}
	data.TotalMapel, err = repository.CountMapel()
	if err != nil {
		return AdminDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat data total mata pelajaran")
	}
	data.TotalPerizinanPending, err = repository.CountPerizinanPending()
	if err != nil {
		return AdminDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat data total perizinan pending")
	}

	return data, nil
}

func GetGuruDashboard(ctx academicModel.TeacherContext) (GuruDashboardData, error) {
	if ctx.Role != "guru" {
		return GuruDashboardData{}, errors.New("ERR_FORBIDDEN: akses hanya untuk guru")
	}

	var data GuruDashboardData
	todayDayName := getIndonesianDayName()
	var err error

	// 1. Today's Schedule Count
	data.TodaySchedule, err = repository.CountTodaySchedule(ctx.GuruID, todayDayName)
	if err != nil {
		return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat jadwal")
	}

	// 2. Total Classes Assigned
	data.TotalClasses, err = repository.CountTotalClasses(ctx.GuruID)
	if err != nil {
		return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat kelas")
	}

	// 3. Total Students in Classes Taught
	data.TotalStudents, err = repository.CountTotalStudentsForGuru(ctx.GuruID)
	if err != nil {
		return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat total siswa")
	}

	// 4. Pending Grade Completion
	schedules, err := repository.GetJadwalByGuru(ctx.GuruID)
	if err != nil {
		return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat jadwal perwalian")
	}

	var pendingGrades int64 = 0
	for _, sched := range schedules {
		classStudentsCount, err := repository.CountSiswaByKelas(sched.KelasID)
		if err != nil {
			return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat data siswa kelas")
		}

		gradedCount, err := repository.CountNilaiByJadwal(sched.MapelID, sched.KelasID, sched.Semester, sched.TahunAjaran)
		if err != nil {
			return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat nilai siswa")
		}

		if classStudentsCount > gradedCount {
			pendingGrades += (classStudentsCount - gradedCount)
		}
	}
	data.PendingGrades = pendingGrades

	// 5. Today's schedules list
	data.JadwalHariIni, err = repository.GetTodayJadwal(ctx.GuruID, todayDayName)
	if err != nil {
		return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat daftar jadwal hari ini")
	}

	// Populate TotalSiswa for each schedule's class
	for i := range data.JadwalHariIni {
		count, err := repository.CountSiswaByKelas(data.JadwalHariIni[i].KelasID)
		if err != nil {
			return GuruDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat total siswa per kelas")
		}
		data.JadwalHariIni[i].Kelas.TotalSiswa = int(count)
	}

	return data, nil
}

func GetSiswaDashboard(userID uint) (SiswaDashboardData, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
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
	data.ProfilSiswa, err = repository.GetSiswaByUserID(user.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return SiswaDashboardData{}, errors.New("ERR_PROFILE_NOT_FOUND: data profil siswa tidak ditemukan")
		}
		return SiswaDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan database")
	}

	// 1. Jadwal Hari Ini untuk kelas siswa tersebut
	data.JadwalHariIni, err = repository.GetTodayJadwalForKelas(data.ProfilSiswa.KelasID, todayDayName)
	if err != nil {
		return SiswaDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat jadwal pelajaran hari ini")
	}

	// 2. Riwayat Perizinan Siswa (Limit 5)
	data.RiwayatPerizinan, err = repository.GetRiwayatPerizinanSiswa(data.ProfilSiswa.ID, 5)
	if err != nil {
		return SiswaDashboardData{}, errors.New("ERR_INTERNAL_SERVER: kegagalan memuat riwayat perizinan")
	}

	return data, nil
}
