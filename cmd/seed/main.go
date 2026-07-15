package main

import (
	"fmt"
	"log"
	"os"

	"backend/config"
	authModel "backend/modules/auth/model"
	guruModel "backend/modules/guru/model"
	jadwalModel "backend/modules/jadwal/model"
	kelasModel "backend/modules/kelas/model"
	mapelModel "backend/modules/mapel/model"
	nilaiModel "backend/modules/nilai/model"
	perizinanModel "backend/modules/perizinan/model"
	absensiModel "backend/modules/absensi/model"
	schoolModel "backend/modules/school/model"
	siswaModel "backend/modules/siswa/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func strPtr(s string) *string {
	return &s
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

type SubjectDef struct {
	Kode         string
	Nama         string
	Jurusan      string
	IsUmum       bool
	TeacherIndex int
	Limit        int
}

// Green/Red console colors
const (
	ColorReset = "\033[0m"
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
)

func printPass(msg string) {
	fmt.Printf("%-20s : %sPASS%s\n", msg, ColorGreen, ColorReset)
}

func printFail(msg string) {
	fmt.Printf("%-20s : %sFAILED%s\n", msg, ColorRed, ColorReset)
}

func main() {
	config.ConnectDB()
	db := config.DB
	log.Println("Starting database seeding process for Release Candidate v2.1...")

	errTransaction := db.Transaction(func(tx *gorm.DB) error {
		log.Println("Clearing existing database records...")
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&nilaiModel.Nilai{})
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&jadwalModel.Jadwal{})
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&guruModel.GuruMapel{})
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&mapelModel.Mapel{})
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&siswaModel.Siswa{})
		tx.Exec("UPDATE kelas SET wali_kelas_id = NULL")
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&guruModel.Guru{})
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&kelasModel.Kelas{})
		// Ignore error if table doesn't exist yet by capturing it in a nested tx or just use GORM
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&perizinanModel.Perizinan{})
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&absensiModel.Absensi{})
		// Continue with gorm deletes
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&authModel.User{})
		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&schoolModel.SchoolProfile{})

		hashAdmin, _ := bcrypt.GenerateFromPassword([]byte("Admin123!"), bcrypt.DefaultCost)
		hashGuru, _ := bcrypt.GenerateFromPassword([]byte("Guru123!"), bcrypt.DefaultCost)
		hashSiswa, _ := bcrypt.GenerateFromPassword([]byte("Siswa123!"), bcrypt.DefaultCost)

		schoolProfile := schoolModel.SchoolProfile{
			ID:                1,
			Name:              "SMK Negeri 1 Salatiga",
			Npsn:              "20312345",
			Status:            "Negeri",
			Level:             "SMK / Sekolah Menengah Kejuruan",
			Accreditation:     "A (Sangat Baik)",
			PrincipalName:     "Drs. Hadi Santoso",
			PrincipalNip:      "GR0001",
			AcademicYear:      "2026/2027",
			CurrentSemester:   "Ganjil",
		}
		if err := tx.Create(&schoolProfile).Error; err != nil {
			return err
		}

		adminUser := authModel.User{
			Nama:         "Administrator",
			Email:        strPtr("admin@sekolah.com"),
			Password:     string(hashAdmin),
			Role:         "admin",
			IsFirstLogin: true,
		}
		if err := tx.Create(&adminUser).Error; err != nil {
			return err
		}

		teacherNames := []string{
			"Drs. Hadi Santoso", "Siti Aminah, M.Pd.", "Joko Susilo, S.Kom.", "Sri Wahyuni, M.E.", "Prabowo Subianto, S.H.",
			"Megawati Sukarno, M.Hum.", "Susilo Bambang, Ph.D.", "Abdurrahman Wahid, Lc.", "Dr. Bacharuddin J. Habibie", "Soeharto Hadi, M.Si.",
			"Agus Harimurti, M.Sc.", "Anies Baswedan, Ph.D.", "Ganjar Pranowo, S.H.", "Ridwan Kamil, S.T.", "Sandiaga Uno, MBA",
			"Erick Thohir, M.B.A.", "Luhut Pandjaitan, M.P.A.", "Basuki Purnama, M.M.", "Tri Rismaharini, M.T.", "Khofifah Indar, M.Si.",
			"Mahfud MD, S.H.", "Muhaimin Iskandar, M.Si.", "Puan Maharani, S.Sos.", "Gibran Rakabuming, B.Sc.", "Djarot Saiful, M.S.",
		}

		teachers := make([]guruModel.Guru, 25)
		for i := 0; i < 25; i++ {
			emailStr := "principal@sekolah.com"
			if i > 0 {
				emailStr = fmt.Sprintf("guru%02d@sekolah.com", i+1)
			}

			u := authModel.User{
				Nama:         teacherNames[i],
				Email:        strPtr(emailStr),
				Password:     string(hashGuru),
				Role:         "guru",
				IsFirstLogin: true,
			}
			if err := tx.Create(&u).Error; err != nil {
				return err
			}

			gender := "Laki-laki"
			if i == 1 || i == 3 || i == 5 || i == 18 || i == 19 || i == 22 {
				gender = "Perempuan"
			}

			g := guruModel.Guru{
				UserID:       u.ID,
				NIP:          fmt.Sprintf("GR%04d", i+1),
				Nama:         teacherNames[i],
				Gelar:        "S.Pd.",
				JenisKelamin: gender,
				NoHP:         fmt.Sprintf("0812345678%02d", i+1),
				Provinsi:     "JAWA TENGAH",
				Kabupaten:    "KOTA SALATIGA",
				Kecamatan:    "SIDOREJO",
				Desa:         "SIDOREJO LOR",
				AlamatDetail: fmt.Sprintf("Jl. Diponegoro No. %d", i+1),
			}
			if err := tx.Create(&g).Error; err != nil {
				return err
			}
			teachers[i] = g
		}

		classNames := []string{"XI-RPL-1", "XI-RPL-2", "XI-TKJ-1", "XI-TKJ-2", "XI-DKV-1"}
		jurusans := []string{"RPL", "RPL", "TKJ", "TKJ", "DKV"}
		classes := make([]kelasModel.Kelas, 5)

		for i := 0; i < 5; i++ {
			k := kelasModel.Kelas{
				NamaKelas:   classNames[i],
				Tingkat:     "11",
				Jurusan:     jurusans[i],
				Kapasitas:   25,
				WaliKelasID: &teachers[i].ID,
			}
			if err := tx.Create(&k).Error; err != nil {
				return err
			}
			classes[i] = k
		}

		firstNames := []string{"Rian", "Budi", "Santi", "Dewi", "Joko", "Megawati", "Susilo", "Abdurrahman", "Habibie", "Soeharto", "Agus", "Anies", "Ganjar", "Ridwan", "Sandiaga", "Erick", "Luhut", "Ahok", "Risma", "Khofifah", "Dimas", "Bagus", "Fajar", "Dwi", "Eko", "Tri", "Yanto", "Rini", "Wulan", "Aditya"}
		lastNames := []string{"Pratama", "Hidayat", "Saputra", "Wibowo", "Kurniawan", "Setiawan", "Nugroho", "Suryono", "Santoso", "Budiman", "Hartono", "Gunawan", "Susanto", "Wijaya", "Siregar"}
		students := make([]siswaModel.Siswa, 125)
		
		for i := 0; i < 125; i++ {
			name := firstNames[i%len(firstNames)] + " " + lastNames[i%len(lastNames)]
			emailStr := fmt.Sprintf("siswa%03d@sekolah.com", i+1)

			u := authModel.User{
				Nama:         name,
				Email:        strPtr(emailStr),
				Password:     string(hashSiswa),
				Role:         "siswa",
				IsFirstLogin: true,
			}
			if err := tx.Create(&u).Error; err != nil {
				return err
			}

			classIdx := i / 25
			gender := "Laki-laki"
			if i%2 == 1 {
				gender = "Perempuan"
			}

			s := siswaModel.Siswa{
				UserID:       u.ID,
				NIS:          fmt.Sprintf("2026%04d", i+1),
				Nama:         name,
				JenisKelamin: gender,
				TanggalLahir: fmt.Sprintf("2009-08-%02d", (i%28)+1),
				KelasID:      classes[classIdx].ID,
				Provinsi:     "JAWA TENGAH",
				Kabupaten:    "KOTA SALATIGA",
				Kecamatan:    "SIDOREJO",
				Desa:         "SIDOREJO LOR",
				AlamatDetail: fmt.Sprintf("Sidorejo Indah RT 02 RW %02d", classIdx+1),
			}
			if err := tx.Create(&s).Error; err != nil {
				return err
			}
			students[i] = s
		}

		subjectList := []struct{ Kode, Nama string; IsUmum bool; Jurusan string }{
			{Kode: "AGM11", Nama: "Pendidikan Agama", IsUmum: true, Jurusan: ""},
			{Kode: "PKN11", Nama: "Pendidikan Pancasila", IsUmum: true, Jurusan: ""},
			{Kode: "IND11", Nama: "Bahasa Indonesia", IsUmum: true, Jurusan: ""},
			{Kode: "MTK11", Nama: "Matematika", IsUmum: true, Jurusan: ""},
			{Kode: "SEJ11", Nama: "Sejarah Indonesia", IsUmum: true, Jurusan: ""},
			{Kode: "ING11", Nama: "Bahasa Inggris", IsUmum: true, Jurusan: ""},
			{Kode: "KWU11", Nama: "Kewirausahaan", IsUmum: true, Jurusan: ""},
			{Kode: "P511", Nama: "Project P5", IsUmum: true, Jurusan: ""},
			{Kode: "RPL11", Nama: "Pemrograman Web Dasar", IsUmum: false, Jurusan: "RPL"},
			{Kode: "RPL12", Nama: "Basis Data Sistem", IsUmum: false, Jurusan: "RPL"},
			{Kode: "RPL13", Nama: "Pemrograman Berorientasi Objek", IsUmum: false, Jurusan: "RPL"},
			{Kode: "RPL14", Nama: "Pemrograman Web Lanjut", IsUmum: false, Jurusan: "RPL"},
			{Kode: "TKJ11", Nama: "Administrasi Infrastruktur Jaringan", IsUmum: false, Jurusan: "TKJ"},
			{Kode: "TKJ12", Nama: "Administrasi Sistem Jaringan", IsUmum: false, Jurusan: "TKJ"},
			{Kode: "TKJ13", Nama: "Teknologi Layanan Jaringan", IsUmum: false, Jurusan: "TKJ"},
			{Kode: "TKJ14", Nama: "Keamanan Jaringan & Sistem", IsUmum: false, Jurusan: "TKJ"},
			{Kode: "DKV11", Nama: "Dasar Desain Grafis Komputer", IsUmum: false, Jurusan: "DKV"},
			{Kode: "DKV12", Nama: "Fotografi & Kamera Studio", IsUmum: false, Jurusan: "DKV"},
			{Kode: "DKV13", Nama: "Videografi & Editing Digital", IsUmum: false, Jurusan: "DKV"},
			{Kode: "DKV14", Nama: "Desain Grafis Percetakan Media", IsUmum: false, Jurusan: "DKV"},
		}

		mapelMap := make(map[string]mapelModel.Mapel)
		for _, s := range subjectList {
			m := mapelModel.Mapel{KodeMapel: s.Kode, NamaMapel: s.Nama, Jam: 3, Jurusan: s.Jurusan, IsUmum: s.IsUmum}
			if err := tx.Create(&m).Error; err != nil {
				return err
			}
			mapelMap[s.Kode] = m
		}

		subjectDefs := []SubjectDef{
			{Kode: "AGM11", IsUmum: true, TeacherIndex: 0, Limit: 1},
			{Kode: "PKN11", IsUmum: true, TeacherIndex: 1, Limit: 1},
			{Kode: "IND11", IsUmum: true, TeacherIndex: 2, Limit: 1},
			{Kode: "MTK11", IsUmum: true, TeacherIndex: 3, Limit: 1},
			{Kode: "SEJ11", IsUmum: true, TeacherIndex: 4, Limit: 2},
			{Kode: "ING11", IsUmum: true, TeacherIndex: 5, Limit: 1},
			{Kode: "KWU11", IsUmum: true, TeacherIndex: 6, Limit: 3},
			{Kode: "P511", IsUmum: true, TeacherIndex: 7, Limit: 3},
			{Kode: "ING11", IsUmum: true, TeacherIndex: 20, Limit: 1},
			{Kode: "MTK11", IsUmum: true, TeacherIndex: 21, Limit: 1},
			{Kode: "IND11", IsUmum: true, TeacherIndex: 22, Limit: 1},
			{Kode: "AGM11", IsUmum: true, TeacherIndex: 23, Limit: 1},
			{Kode: "PKN11", IsUmum: true, TeacherIndex: 24, Limit: 1},
			{Kode: "RPL11", Jurusan: "RPL", TeacherIndex: 8, Limit: 3},
			{Kode: "RPL12", Jurusan: "RPL", TeacherIndex: 9, Limit: 3},
			{Kode: "RPL13", Jurusan: "RPL", TeacherIndex: 10, Limit: 3},
			{Kode: "RPL14", Jurusan: "RPL", TeacherIndex: 11, Limit: 3},
			{Kode: "TKJ11", Jurusan: "TKJ", TeacherIndex: 12, Limit: 3},
			{Kode: "TKJ12", Jurusan: "TKJ", TeacherIndex: 13, Limit: 3},
			{Kode: "TKJ13", Jurusan: "TKJ", TeacherIndex: 14, Limit: 3},
			{Kode: "TKJ14", Jurusan: "TKJ", TeacherIndex: 15, Limit: 3},
			{Kode: "DKV11", Jurusan: "DKV", TeacherIndex: 16, Limit: 3},
			{Kode: "DKV12", Jurusan: "DKV", TeacherIndex: 17, Limit: 3},
			{Kode: "DKV13", Jurusan: "DKV", TeacherIndex: 18, Limit: 3},
			{Kode: "DKV14", Jurusan: "DKV", TeacherIndex: 19, Limit: 3},
		}

		for _, def := range subjectDefs {
			mObj := mapelMap[def.Kode]
			gm := guruModel.GuruMapel{GuruID: teachers[def.TeacherIndex].ID, MapelID: mObj.ID}
			tx.Create(&gm)
		}

		var grid [5][5][6]int
		var teacherBusy [25][5][6]bool
		var subjectUsage [5][25]int
		for c := 0; c < 5; c++ {
			for d := 0; d < 5; d++ {
				for s := 0; s < 6; s++ {
					grid[c][d][s] = -1
				}
			}
		}

		var assign func(day, slot, classIdx int) bool
		assign = func(day, slot, classIdx int) bool {
			if classIdx == 5 {
				nextSlot, nextDay := slot+1, day
				if nextSlot == 6 {
					nextSlot, nextDay = 0, day+1
				}
				if nextDay == 5 {
					return true
				}
				return assign(nextDay, nextSlot, 0)
			}
			class := classes[classIdx]
			for subIdx, def := range subjectDefs {
				if !def.IsUmum && def.Jurusan != class.Jurusan {
					continue
				}
				if subjectUsage[classIdx][subIdx] >= def.Limit {
					continue
				}
				tIdx := def.TeacherIndex
				if teacherBusy[tIdx][day][slot] {
					continue
				}
				grid[classIdx][day][slot] = subIdx
				teacherBusy[tIdx][day][slot] = true
				subjectUsage[classIdx][subIdx]++
				if assign(day, slot, classIdx+1) {
					return true
				}
				grid[classIdx][day][slot] = -1
				teacherBusy[tIdx][day][slot] = false
				subjectUsage[classIdx][subIdx]--
			}
			return false
		}

		if !assign(0, 0, 0) {
			return fmt.Errorf("failed to solve schedule without conflicts")
		}

		days := []string{"Senin", "Selasa", "Rabu", "Kamis", "Jumat"}
		jamMulais := []string{"07:00:00", "08:00:00", "09:00:00", "10:30:00", "11:30:00", "13:30:00"}
		jamSelesaise := []string{"08:00:00", "09:00:00", "10:00:00", "11:30:00", "12:30:00", "14:30:00"}

		for cIdx := 0; cIdx < 5; cIdx++ {
			for dIdx := 0; dIdx < 5; dIdx++ {
				for sIdx := 0; sIdx < 6; sIdx++ {
					subIdx := grid[cIdx][dIdx][sIdx]
					if subIdx == -1 {
						continue
					}
					sched := jadwalModel.Jadwal{
						KelasID:     classes[cIdx].ID,
						MapelID:     mapelMap[subjectDefs[subIdx].Kode].ID,
						GuruID:      teachers[subjectDefs[subIdx].TeacherIndex].ID,
						Hari:        days[dIdx],
						JamMulai:    jamMulais[sIdx],
						JamSelesai:  jamSelesaise[sIdx],
						TahunAjaran: "2026/2027",
						Semester:    "Ganjil",
					}
					if err := tx.Create(&sched).Error; err != nil {
						return err
					}
				}
			}
		}

		var gradesList []nilaiModel.Nilai
		for _, s := range students {
			var scheds []jadwalModel.Jadwal
			tx.Where("kelas_id = ? AND tahun_ajaran = ?", s.KelasID, "2026/2027").Find(&scheds)

			uniqueScheds := make(map[uint]jadwalModel.Jadwal)
			for _, sc := range scheds {
				uniqueScheds[sc.MapelID] = sc
			}

			for _, semester := range []string{"Ganjil", "Genap"} {
				semVal := 1
				if semester == "Genap" {
					semVal = 2
				}
				for mapelID, sc := range uniqueScheds {
					studentID := s.ID
					// Realistic Grade Distribution (Pseudorandom but deterministic)
					baseHash := (studentID * 37) + (mapelID * 17) + (uint(semVal) * 11)
					
					// Tugas: 70-100
					tugas := 70.0 + float64(baseHash%31)
					// UTS: 65-100
					uts := 65.0 + float64((baseHash*2)%36)
					// UAS: 60-100
					uas := 60.0 + float64((baseHash*3)%41)

					// Adjust to target A(20%), B(40%), C(30%), D(10%)
					distVal := baseHash % 100
					if distVal < 20 { // A (>=85)
						tugas = 90.0 + float64(baseHash%11)
						uts = 85.0 + float64((baseHash*2)%16)
						uas = 85.0 + float64((baseHash*3)%16)
					} else if distVal < 60 { // B (75-84)
						tugas = 80.0 + float64(baseHash%10)
						uts = 75.0 + float64((baseHash*2)%10)
						uas = 75.0 + float64((baseHash*3)%10)
					} else if distVal < 90 { // C (65-74)
						tugas = 70.0 + float64(baseHash%10)
						uts = 65.0 + float64((baseHash*2)%10)
						uas = 65.0 + float64((baseHash*3)%10)
					} else { // D (<65)
						tugas = 60.0 + float64(baseHash%10)
						uts = 60.0 + float64((baseHash*2)%10)
						uas = 60.0 + float64((baseHash*3)%10)
					}

					finalScore := (tugas * 0.3) + (uts * 0.3) + (uas * 0.4)
					gradeLetter := deriveGradeHuruf(finalScore)

					gradesList = append(gradesList, nilaiModel.Nilai{
						SiswaID:     studentID,
						KelasID:     s.KelasID,
						MapelID:     mapelID,
						GuruID:      sc.GuruID,
						Tugas:       tugas,
						UTS:         uts,
						UAS:         uas,
						NilaiAkhir:  finalScore,
						GradeHuruf:  gradeLetter,
						Semester:    semester,
						TahunAjaran: "2026/2027",
					})
				}
			}
		}

		if err := tx.CreateInBatches(gradesList, 100).Error; err != nil {
			return err
		}

		return verifySeededDataAndReport(tx)
	})

	if errTransaction != nil {
		fmt.Printf("\n%sSeeder encountered an error and rolled back: %v%s\n", ColorRed, errTransaction, ColorReset)
		os.Exit(1)
	}

	os.Exit(0)
}

func verifySeededDataAndReport(tx *gorm.DB) error {
	fmt.Println("\n================================")
	fmt.Println("DATABASE HEALTH REPORT")
	fmt.Println("================================")

	var uCount, gCount, sCount, kCount, mCount, gmCount, jCount, nCount, pCount int64
	tx.Model(&authModel.User{}).Count(&uCount)
	tx.Model(&guruModel.Guru{}).Count(&gCount)
	tx.Model(&siswaModel.Siswa{}).Count(&sCount)
	tx.Model(&kelasModel.Kelas{}).Count(&kCount)
	tx.Model(&mapelModel.Mapel{}).Count(&mCount)
	tx.Model(&guruModel.GuruMapel{}).Count(&gmCount)
	tx.Model(&jadwalModel.Jadwal{}).Count(&jCount)
	tx.Model(&nilaiModel.Nilai{}).Count(&nCount)
	tx.Model(&schoolModel.SchoolProfile{}).Count(&pCount)

	fmt.Printf("%-20s : %d\n", "Users", uCount)
	fmt.Printf("%-20s : %d\n", "Guru", gCount)
	fmt.Printf("%-20s : %d\n", "Siswa", sCount)
	fmt.Printf("%-20s : %d\n", "Kelas", kCount)
	fmt.Printf("%-20s : %d\n", "Mapel", mCount)
	fmt.Printf("%-20s : %d\n", "GuruMapel", gmCount)
	fmt.Printf("%-20s : %d\n", "Jadwal", jCount)
	fmt.Printf("%-20s : %d\n", "Nilai", nCount)
	fmt.Println("================================")

	allPass := true
	var errs []string

	// Integrity checks
	var userGuruMismatch int64
	tx.Model(&guruModel.Guru{}).Where("user_id IS NULL OR user_id NOT IN (SELECT id FROM users)").Count(&userGuruMismatch)
	if userGuruMismatch > 0 { errs = append(errs, "Guru without User") }

	var gmGuruMismatch int64
	tx.Model(&guruModel.GuruMapel{}).Where("guru_id NOT IN (SELECT id FROM gurus)").Count(&gmGuruMismatch)
	if gmGuruMismatch > 0 { errs = append(errs, "GuruMapel without Guru") }

	var schedWithoutGuruMapel int64
	tx.Raw(`SELECT count(*) FROM jadwals j LEFT JOIN guru_mapel gm ON j.guru_id = gm.guru_id AND j.mapel_id = gm.mapel_id WHERE gm.id IS NULL`).Scan(&schedWithoutGuruMapel)
	if schedWithoutGuruMapel > 0 { errs = append(errs, "Jadwal without matching GuruMapel") }

	var schedBroken int64
	tx.Model(&jadwalModel.Jadwal{}).Where("guru_id NOT IN (SELECT id FROM gurus) OR mapel_id NOT IN (SELECT id FROM mapels) OR kelas_id NOT IN (SELECT id FROM kelas)").Count(&schedBroken)
	if schedBroken > 0 { errs = append(errs, "Jadwal with broken FK") }

	var classSiswaMismatch int64
	tx.Model(&siswaModel.Siswa{}).Where("kelas_id NOT IN (SELECT id FROM kelas)").Count(&classSiswaMismatch)
	if classSiswaMismatch > 0 { errs = append(errs, "Siswa without valid Kelas") }

	var nilaiBroken int64
	tx.Model(&nilaiModel.Nilai{}).Where("guru_id NOT IN (SELECT id FROM gurus) OR mapel_id NOT IN (SELECT id FROM mapels) OR siswa_id NOT IN (SELECT id FROM siswas)").Count(&nilaiBroken)
	if nilaiBroken > 0 { errs = append(errs, "Nilai with broken FK") }

	var duplicateNIS int64
	tx.Raw(`SELECT count(nis) FROM siswas GROUP BY nis HAVING count(nis) > 1`).Count(&duplicateNIS)
	if duplicateNIS > 0 { errs = append(errs, "Duplicate NIS") }

	var duplicateEmail int64
	tx.Raw(`SELECT count(email) FROM users GROUP BY email HAVING count(email) > 1`).Count(&duplicateEmail)
	if duplicateEmail > 0 { errs = append(errs, "Duplicate Email") }

	if len(errs) > 0 {
		printFail("Orphan Data / Broken FK")
		for _, e := range errs {
			fmt.Printf(" - %s\n", e)
		}
		allPass = false
	} else {
		printPass("Orphan Data / Broken FK")
	}

	fmt.Println("================================")
	if allPass {
		fmt.Printf("Overall Status       : %sPASS%s\n", ColorGreen, ColorReset)
		fmt.Printf("Exit Code            : %s0%s\n", ColorGreen, ColorReset)
		fmt.Println("================================")
		return nil
	} else {
		fmt.Printf("Overall Status       : %sFAILED%s\n", ColorRed, ColorReset)
		fmt.Printf("Exit Code            : %s1%s\n", ColorRed, ColorReset)
		fmt.Println("================================")
		return fmt.Errorf("integrity checks failed")
	}
}
