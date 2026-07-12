package main

import (
	"fmt"
	"log"

	"backend/config"
	authModel "backend/modules/auth/model"
	guruModel "backend/modules/guru/model"
	jadwalModel "backend/modules/jadwal/model"
	kelasModel "backend/modules/kelas/model"
	mapelModel "backend/modules/mapel/model"
	nilaiModel "backend/modules/nilai/model"
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

func main() {
	// 1. Connect to Database
	config.ConnectDB()
	db := config.DB
	log.Println("Starting database seeding process for Release Candidate v1.0...")

	// 2. Clear Tables in Idempotent Order (Referential Integrity)
	log.Println("Clearing existing database records for idempotency...")
	db.Exec("DELETE FROM nilais")
	db.Exec("DELETE FROM jadwals")
	db.Exec("DELETE FROM guru_mapel")
	db.Exec("DELETE FROM mapels")
	db.Exec("DELETE FROM siswas")
	db.Exec("UPDATE kelas SET wali_kelas_id = NULL")
	db.Exec("DELETE FROM gurus")
	db.Exec("DELETE FROM kelas")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM school_profiles")
	log.Println("Database tables cleared successfully.")

	// 3. Pre-calculate Bcrypt Hashes to speed up seeding
	log.Println("Pre-calculating bcrypt password hashes...")
	hashAdmin, _ := bcrypt.GenerateFromPassword([]byte("Admin123!"), bcrypt.DefaultCost)
	hashGuru, _ := bcrypt.GenerateFromPassword([]byte("Guru123!"), bcrypt.DefaultCost)
	hashSiswa, _ := bcrypt.GenerateFromPassword([]byte("Siswa123!"), bcrypt.DefaultCost)

	// 4. Seed School Profile Record with Principal NIG GR0001
	log.Println("Seeding School Profile into database...")
	schoolProfile := schoolModel.SchoolProfile{
		ID:                1,
		Logo:              "",
		Name:              "SMK Negeri 1 Salatiga",
		Npsn:              "20312345",
		Status:            "Negeri",
		Level:             "SMK / Sekolah Menengah Kejuruan",
		Accreditation:     "A (Sangat Baik)",
		EstablishedYear:   "1965",
		Address:           "Jl. Diponegoro No. 25, Salatiga",
		PostalCode:        "50711",
		Phone:             "(0298) 123456",
		Email:             "info@smkn1salatiga.sch.id",
		Website:           "https://smkn1salatiga.sch.id",
		PrincipalName:     "Drs. Hadi Santoso",
		PrincipalNip:      "GR0001",
		PrincipalPosition: "Kepala Sekolah",
		AppointmentPeriod: "2021 - 2027",
		AcademicYear:      "2026/2027",
		CurrentSemester:   "Ganjil",
		AcademicStatus:    "🟢 Aktif",
		SchoolType:        "Sekolah Menengah Kejuruan (SMK)",
		Curriculum:        "Kurikulum Merdeka",
		Shift:             "Pagi",
		OperationalStatus: "Aktif",
		Vision:            "Menjadi lembaga pendidikan kejuruan yang unggul, berkarakter, dan berdaya saing global di era digital.",
		Mission:           "Menyelenggarakan pembelajaran berbasis kompetensi teknologi industri dan keahlian terkini.\nMembina karakter siswa berlandaskan iman, taqwa, nilai Pancasila, dan budi pekerti luhur.\nMenjalin kemitraan erat dengan dunia usaha, dunia industri, dan asosiasi profesi (DUDI).\nMengembangkan semangat kewirausahaan, daya saing, dan kemandirian siswa.",
		Facilities:        "Library, Computer Laboratory, Mosque, Workshop, Sports Field",
	}
	if err := db.Create(&schoolProfile).Error; err != nil {
		log.Fatal("Failed seeding School Profile:", err)
	}
	log.Println("School Profile seeded successfully.")

	// 5. Seed Administrator User (Enforced force change password: IsFirstLogin = true)
	adminUser := authModel.User{
		Nama:         "Administrator",
		Email:        strPtr("admin@sekolah.com"),
		Password:     string(hashAdmin),
		Role:         "admin",
		IsFirstLogin: true,
	}
	if err := db.Create(&adminUser).Error; err != nil {
		log.Fatal("Failed seeding Administrator:", err)
	}
	log.Println("Administrator user seeded successfully.")

	// 6. Seed Teachers (exactly 25, starting with Principal as Teacher #0 - GR0001)
	teacherNames := []string{
		"Drs. Hadi Santoso", "Siti Aminah, M.Pd.", "Joko Susilo, S.Kom.", "Sri Wahyuni, M.E.", "Prabowo Subianto, S.H.",
		"Megawati Sukarno, M.Hum.", "Susilo Bambang, Ph.D.", "Abdurrahman Wahid, Lc.", "Dr. Bacharuddin J. Habibie", "Soeharto Hadi, M.Si.",
		"Agus Harimurti, M.Sc.", "Anies Baswedan, Ph.D.", "Ganjar Pranowo, S.H.", "Ridwan Kamil, S.T.", "Sandiaga Uno, MBA",
		"Erick Thohir, M.B.A.", "Luhut Pandjaitan, M.P.A.", "Basuki Purnama, M.M.", "Tri Rismaharini, M.T.", "Khofifah Indar, M.Si.",
		"Mahfud MD, S.H.", "Muhaimin Iskandar, M.Si.", "Puan Maharani, S.Sos.", "Gibran Rakabuming, B.Sc.", "Djarot Saiful, M.S.",
	}

	teachers := make([]guruModel.Guru, 25)
	for i := 0; i < 25; i++ {
		var emailStr string
		if i == 0 {
			emailStr = "principal@sekolah.com"
		} else {
			emailStr = fmt.Sprintf("guru%02d@sekolah.com", i+1)
		}

		u := authModel.User{
			Nama:         teacherNames[i],
			Email:        strPtr(emailStr),
			Password:     string(hashGuru),
			Role:         "guru",
			IsFirstLogin: true,
		}
		if err := db.Create(&u).Error; err != nil {
			log.Fatal("Failed seeding Teacher User:", emailStr, err)
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
		if err := db.Create(&g).Error; err != nil {
			log.Fatal("Failed seeding Teacher profile:", g.Nama, err)
		}
		teachers[i] = g
	}
	log.Println("Seeded 25 Teachers successfully.")

	// 7. Seed exactly 5 Classes
	classNames := []string{"XI-RPL-1", "XI-RPL-2", "XI-TKJ-1", "XI-TKJ-2", "XI-DKV-1"}
	jurusans := []string{"RPL", "RPL", "TKJ", "TKJ", "DKV"}
	classes := make([]kelasModel.Kelas, 5)

	for i := 0; i < 5; i++ {
		k := kelasModel.Kelas{
			NamaKelas:   classNames[i],
			Tingkat:     "11",
			Jurusan:     jurusans[i],
			Kapasitas:   25,
			WaliKelasID: &teachers[i].ID, // homeroom teacher assignments (first 5 teachers)
		}
		if err := db.Create(&k).Error; err != nil {
			log.Fatal("Failed seeding Class:", classNames[i], err)
		}
		classes[i] = k
	}
	log.Println("Seeded 5 Classes successfully.")

	// 8. Seed exactly 125 Students (exactly 25 per class, with NIS 20260001 - 20260125)
	firstNames := []string{
		"Rian", "Budi", "Santi", "Dewi", "Joko", "Megawati", "Susilo", "Abdurrahman", "Habibie", "Soeharto",
		"Agus", "Anies", "Ganjar", "Ridwan", "Sandiaga", "Erick", "Luhut", "Ahok", "Risma", "Khofifah",
		"Dimas", "Bagus", "Fajar", "Dwi", "Eko", "Tri", "Yanto", "Rini", "Wulan", "Aditya",
	}
	lastNames := []string{
		"Pratama", "Hidayat", "Saputra", "Wibowo", "Kurniawan", "Setiawan", "Nugroho", "Suryono", "Santoso", "Budiman",
		"Hartono", "Gunawan", "Susanto", "Wijaya", "Siregar",
	}

	students := make([]siswaModel.Siswa, 125)
	for i := 0; i < 125; i++ {
		first := firstNames[i%len(firstNames)]
		last := lastNames[i%len(lastNames)]
		name := first + " " + last
		emailStr := fmt.Sprintf("siswa%03d@sekolah.com", i+1)

		u := authModel.User{
			Nama:         name,
			Email:        strPtr(emailStr),
			Password:     string(hashSiswa),
			Role:         "siswa",
			IsFirstLogin: true,
		}
		if err := db.Create(&u).Error; err != nil {
			log.Fatal("Failed seeding Student User:", emailStr, err)
		}

		classIdx := i / 25 // exactly 25 students per class
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
			NoHP:         fmt.Sprintf("082134567%03d", i+1),
			KelasID:      classes[classIdx].ID,
			Provinsi:     "JAWA TENGAH",
			Kabupaten:    "KOTA SALATIGA",
			Kecamatan:    "SIDOREJO",
			Desa:         "SIDOREJO LOR",
			AlamatDetail: fmt.Sprintf("Sidorejo Indah RT 02 RW %02d", classIdx+1),
		}
		if err := db.Create(&s).Error; err != nil {
			log.Fatal("Failed seeding Student Profile:", name, err)
		}
		students[i] = s
	}
	log.Println("Seeded 125 Students (exactly 25 per class) successfully.")

	// 9. Seed exactly 20 Subjects
	subjectList := []struct {
		Kode    string
		Nama    string
		IsUmum  bool
		Jurusan string
	}{
		// 8 General Subjects
		{Kode: "AGM11", Nama: "Pendidikan Agama", IsUmum: true, Jurusan: ""},
		{Kode: "PKN11", Nama: "Pendidikan Pancasila", IsUmum: true, Jurusan: ""},
		{Kode: "IND11", Nama: "Bahasa Indonesia", IsUmum: true, Jurusan: ""},
		{Kode: "MTK11", Nama: "Matematika", IsUmum: true, Jurusan: ""},
		{Kode: "SEJ11", Nama: "Sejarah Indonesia", IsUmum: true, Jurusan: ""},
		{Kode: "ING11", Nama: "Bahasa Inggris", IsUmum: true, Jurusan: ""},
		{Kode: "KWU11", Nama: "Kewirausahaan", IsUmum: true, Jurusan: ""},
		{Kode: "P511", Nama: "Project P5", IsUmum: true, Jurusan: ""},

		// 4 RPL Subjects
		{Kode: "RPL11", Nama: "Pemrograman Web Dasar", IsUmum: false, Jurusan: "RPL"},
		{Kode: "RPL12", Nama: "Basis Data Sistem", IsUmum: false, Jurusan: "RPL"},
		{Kode: "RPL13", Nama: "Pemrograman Berorientasi Objek", IsUmum: false, Jurusan: "RPL"},
		{Kode: "RPL14", Nama: "Pemrograman Web Lanjut", IsUmum: false, Jurusan: "RPL"},

		// 4 TKJ Subjects
		{Kode: "TKJ11", Nama: "Administrasi Infrastruktur Jaringan", IsUmum: false, Jurusan: "TKJ"},
		{Kode: "TKJ12", Nama: "Administrasi Sistem Jaringan", IsUmum: false, Jurusan: "TKJ"},
		{Kode: "TKJ13", Nama: "Teknologi Layanan Jaringan", IsUmum: false, Jurusan: "TKJ"},
		{Kode: "TKJ14", Nama: "Keamanan Jaringan & Sistem", IsUmum: false, Jurusan: "TKJ"},

		// 4 DKV Subjects
		{Kode: "DKV11", Nama: "Dasar Desain Grafis Komputer", IsUmum: false, Jurusan: "DKV"},
		{Kode: "DKV12", Nama: "Fotografi & Kamera Studio", IsUmum: false, Jurusan: "DKV"},
		{Kode: "DKV13", Nama: "Videografi & Editing Digital", IsUmum: false, Jurusan: "DKV"},
		{Kode: "DKV14", Nama: "Desain Grafis Percetakan Media", IsUmum: false, Jurusan: "DKV"},
	}

	subjects := make([]mapelModel.Mapel, len(subjectList))
	mapelMap := make(map[string]mapelModel.Mapel)

	for idx, s := range subjectList {
		m := mapelModel.Mapel{
			KodeMapel: s.Kode,
			NamaMapel: s.Nama,
			Jam:       3,
			Jurusan:   s.Jurusan,
			IsUmum:    s.IsUmum,
		}
		if err := db.Create(&m).Error; err != nil {
			log.Fatal("Failed seeding Subject:", s.Nama, err)
		}
		subjects[idx] = m
		mapelMap[s.Kode] = m
	}
	log.Println("Seeded exactly 20 subjects.")

	// Define 25 Subject solver definitions (mapping 20 subjects to 25 teachers)
	subjectDefs := []SubjectDef{
		// General Subjects (mapped to Teachers 0-7)
		{Kode: "AGM11", Nama: "Pendidikan Agama", Jurusan: "", IsUmum: true, TeacherIndex: 0, Limit: 1},
		{Kode: "PKN11", Nama: "Pendidikan Pancasila", Jurusan: "", IsUmum: true, TeacherIndex: 1, Limit: 1},
		{Kode: "IND11", Nama: "Bahasa Indonesia", Jurusan: "", IsUmum: true, TeacherIndex: 2, Limit: 1},
		{Kode: "MTK11", Nama: "Matematika", Jurusan: "", IsUmum: true, TeacherIndex: 3, Limit: 1},
		{Kode: "SEJ11", Nama: "Sejarah Indonesia", Jurusan: "", IsUmum: true, TeacherIndex: 4, Limit: 2},
		{Kode: "ING11", Nama: "Bahasa Inggris", Jurusan: "", IsUmum: true, TeacherIndex: 5, Limit: 1},
		{Kode: "KWU11", Nama: "Kewirausahaan", Jurusan: "", IsUmum: true, TeacherIndex: 6, Limit: 3},
		{Kode: "P511", Nama: "Project P5", Jurusan: "", IsUmum: true, TeacherIndex: 7, Limit: 3},

		// Helper/assistant teachers for general subjects (mapped to Teachers 20-24)
		{Kode: "ING11", Nama: "Bahasa Inggris", Jurusan: "", IsUmum: true, TeacherIndex: 20, Limit: 1},
		{Kode: "MTK11", Nama: "Matematika", Jurusan: "", IsUmum: true, TeacherIndex: 21, Limit: 1},
		{Kode: "IND11", Nama: "Bahasa Indonesia", Jurusan: "", IsUmum: true, TeacherIndex: 22, Limit: 1},
		{Kode: "AGM11", Nama: "Pendidikan Agama", Jurusan: "", IsUmum: true, TeacherIndex: 23, Limit: 1},
		{Kode: "PKN11", Nama: "Pendidikan Pancasila", Jurusan: "", IsUmum: true, TeacherIndex: 24, Limit: 1},

		// RPL Vocational Subjects (mapped to Teachers 8-11)
		{Kode: "RPL11", Nama: "Pemrograman Web Dasar", Jurusan: "RPL", IsUmum: false, TeacherIndex: 8, Limit: 3},
		{Kode: "RPL12", Nama: "Basis Data Sistem", Jurusan: "RPL", IsUmum: false, TeacherIndex: 9, Limit: 3},
		{Kode: "RPL13", Nama: "Pemrograman Berorientasi Objek", Jurusan: "RPL", IsUmum: false, TeacherIndex: 10, Limit: 3},
		{Kode: "RPL14", Nama: "Pemrograman Web Lanjut", Jurusan: "RPL", IsUmum: false, TeacherIndex: 11, Limit: 3},

		// TKJ Vocational Subjects (mapped to Teachers 12-15)
		{Kode: "TKJ11", Nama: "Administrasi Infrastruktur Jaringan", Jurusan: "TKJ", IsUmum: false, TeacherIndex: 12, Limit: 3},
		{Kode: "TKJ12", Nama: "Administrasi Sistem Jaringan", Jurusan: "TKJ", IsUmum: false, TeacherIndex: 13, Limit: 3},
		{Kode: "TKJ13", Nama: "Teknologi Layanan Jaringan", Jurusan: "TKJ", IsUmum: false, TeacherIndex: 14, Limit: 3},
		{Kode: "TKJ14", Nama: "Keamanan Jaringan & Sistem", Jurusan: "TKJ", IsUmum: false, TeacherIndex: 15, Limit: 3},

		// DKV Vocational Subjects (mapped to Teachers 16-19)
		{Kode: "DKV11", Nama: "Dasar Desain Grafis Komputer", Jurusan: "DKV", IsUmum: false, TeacherIndex: 16, Limit: 3},
		{Kode: "DKV12", Nama: "Fotografi & Kamera Studio", Jurusan: "DKV", IsUmum: false, TeacherIndex: 17, Limit: 3},
		{Kode: "DKV13", Nama: "Videografi & Editing Digital", Jurusan: "DKV", IsUmum: false, TeacherIndex: 18, Limit: 3},
		{Kode: "DKV14", Nama: "Desain Grafis Percetakan Media", Jurusan: "DKV", IsUmum: false, TeacherIndex: 19, Limit: 3},
	}

	// Register all 25 associations in guru_mapel
	for _, def := range subjectDefs {
		mapelObj := mapelMap[def.Kode]
		gm := guruModel.GuruMapel{
			GuruID:  teachers[def.TeacherIndex].ID,
			MapelID: mapelObj.ID,
		}
		db.Create(&gm)
	}

	// 10. Timetable solver generating conflict-free teaching schedules (exactly 150)
	var grid [5][5][6]int
	var teacherBusy [25][5][6]bool
	var subjectUsage [5][25]int

	// Initialize schedule grid with -1 (empty)
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
			nextSlot := slot + 1
			nextDay := day
			if nextSlot == 6 {
				nextSlot = 0
				nextDay = day + 1
			}
			if nextDay == 5 {
				return true
			}
			return assign(nextDay, nextSlot, 0)
		}

		class := classes[classIdx]
		major := class.Jurusan

		for subIdx, def := range subjectDefs {
			// Validate subject major
			if !def.IsUmum && def.Jurusan != major {
				continue
			}

			// Validate slot limits per class
			if subjectUsage[classIdx][subIdx] >= def.Limit {
				continue
			}

			tIdx := def.TeacherIndex
			if teacherBusy[tIdx][day][slot] {
				continue
			}

			// Make the assignment
			grid[classIdx][day][slot] = subIdx
			teacherBusy[tIdx][day][slot] = true
			subjectUsage[classIdx][subIdx]++

			if assign(day, slot, classIdx+1) {
				return true
			}

			// Backtrack
			grid[classIdx][day][slot] = -1
			teacherBusy[tIdx][day][slot] = false
			subjectUsage[classIdx][subIdx]--
		}

		return false
	}

	log.Println("Solving conflict-free 150 timetable slots...")
	if !assign(0, 0, 0) {
		log.Fatal("Failed to solve schedule timetable without conflicts!")
	}
	log.Println("Timetable grid solved successfully. Saving schedules...")

	days := []string{"Senin", "Selasa", "Rabu", "Kamis", "Jumat"}
	jamMulais := []string{"07:00:00", "08:00:00", "09:00:00", "10:30:00", "11:30:00", "13:30:00"}
	jamSelesaise := []string{"08:00:00", "09:00:00", "10:00:00", "11:30:00", "12:30:00", "14:30:00"}

	schedulesCount := 0
	for cIdx := 0; cIdx < 5; cIdx++ {
		for dIdx := 0; dIdx < 5; dIdx++ {
			for sIdx := 0; sIdx < 6; sIdx++ {
				subIdx := grid[cIdx][dIdx][sIdx]
				if subIdx == -1 {
					continue
				}

				mapelObj := mapelMap[subjectDefs[subIdx].Kode]
				sched := jadwalModel.Jadwal{
					KelasID:     classes[cIdx].ID,
					MapelID:     mapelObj.ID,
					GuruID:      teachers[subjectDefs[subIdx].TeacherIndex].ID,
					Hari:        days[dIdx],
					JamMulai:    jamMulais[sIdx],
					JamSelesai:  jamSelesaise[sIdx],
					TahunAjaran: "2026/2027",
					Semester:    "Ganjil", // schedules are based on standard curriculum
				}
				if err := db.Create(&sched).Error; err != nil {
					log.Fatal("Failed seeding Schedule entry:", err)
				}
				schedulesCount++
			}
		}
	}
	log.Printf("Seeded exactly %d conflict-free teaching schedules.\n", schedulesCount)

	// 11. Seed Academic Grades for Ganjil & Genap Semesters (125 Students * 12 Subjects * 2 Semesters = 3000 Grades)
	log.Println("Bulk seeding 3000 grade records...")
	var gradesList []nilaiModel.Nilai

	for _, s := range students {
		// Get unique subjects assigned to the class schedules
		var studentSchedules []jadwalModel.Jadwal
		db.Where("kelas_id = ? AND tahun_ajaran = ?", s.KelasID, "2026/2027").Find(&studentSchedules)

		uniqueScheds := make(map[uint]jadwalModel.Jadwal)
		for _, sc := range studentSchedules {
			uniqueScheds[sc.MapelID] = sc
		}

		for _, semester := range []string{"Ganjil", "Genap"} {
			semVal := 1
			if semester == "Genap" {
				semVal = 2
			}
			for mapelID, sc := range uniqueScheds {
				studentID := s.ID

				tugas := 70.0 + float64((studentID*3+mapelID*5+uint(semVal)*7)%29)
				uts := 70.0 + float64((studentID*7+mapelID*2+uint(semVal)*13)%29)
				uas := 70.0 + float64((studentID*11+mapelID*13+uint(semVal)*17)%29)

				finalScore := (tugas * 0.3) + (uts * 0.3) + (uas * 0.4)
				gradeLetter := deriveGradeHuruf(finalScore)

				grade := nilaiModel.Nilai{
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
				}
				gradesList = append(gradesList, grade)
			}
		}
	}

	log.Printf("Inserting %d academic grade records...\n", len(gradesList))
	if err := db.CreateInBatches(gradesList, 100).Error; err != nil {
		log.Fatal("Failed bulk seeding Grade records:", err)
	}
	log.Printf("Seeded exactly %d academic grade records.\n", len(gradesList))

	// 12. Run Seeder Verification Assertions
	verifySeededData(db)

	log.Println("Seeding process completed successfully!")
}

func verifySeededData(db *gorm.DB) {
	log.Println("Running automated database verification checks...")

	// 1. Every class contains exactly 25 students.
	var classes []kelasModel.Kelas
	if err := db.Find(&classes).Error; err != nil {
		log.Fatal("Verification failed: cannot query classes:", err)
	}
	for _, c := range classes {
		var studentCount int64
		db.Model(&siswaModel.Siswa{}).Where("kelas_id = ?", c.ID).Count(&studentCount)
		if studentCount != 25 {
			log.Fatalf("Verification failed: class %s has %d students (expected exactly 25)\n", c.NamaKelas, studentCount)
		}
	}
	log.Println("✓ Check 1 passed: Every class contains exactly 25 students.")

	// 2. Every class has one homeroom teacher.
	for _, c := range classes {
		if c.WaliKelasID == nil || *c.WaliKelasID == 0 {
			log.Fatalf("Verification failed: class %s has no homeroom teacher\n", c.NamaKelas)
		}
	}
	log.Println("✓ Check 2 passed: Every class has one homeroom teacher.")

	// 3. Every teacher has at least one teaching schedule.
	var teachers []guruModel.Guru
	if err := db.Find(&teachers).Error; err != nil {
		log.Fatal("Verification failed: cannot query teachers:", err)
	}
	for _, t := range teachers {
		var scheduleCount int64
		db.Model(&jadwalModel.Jadwal{}).Where("guru_id = ?", t.ID).Count(&scheduleCount)
		if scheduleCount == 0 {
			log.Fatalf("Verification failed: teacher %s (ID %d) has 0 teaching schedules\n", t.Nama, t.ID)
		}
	}
	log.Println("✓ Check 3 passed: Every teacher has at least one teaching schedule.")

	// 4. Every subject is assigned to at least one teacher.
	var subjects []mapelModel.Mapel
	if err := db.Find(&subjects).Error; err != nil {
		log.Fatal("Verification failed: cannot query subjects:", err)
	}
	if len(subjects) != 20 {
		log.Fatalf("Verification failed: total subject count is %d (expected exactly 20)\n", len(subjects))
	}
	for _, s := range subjects {
		var relationCount int64
		db.Table("guru_mapel").Where("mapel_id = ?", s.ID).Count(&relationCount)
		if relationCount == 0 {
			log.Fatalf("Verification failed: subject %s (ID %d) is not assigned to any teacher in guru_mapel\n", s.NamaMapel, s.ID)
		}
	}
	log.Println("✓ Check 4 passed: Exactly 20 subjects seeded and assigned to teachers.")

	// 5. Every student has grades.
	var students []siswaModel.Siswa
	if err := db.Find(&students).Error; err != nil {
		log.Fatal("Verification failed: cannot query students:", err)
	}
	for _, st := range students {
		var gradeCount int64
		db.Model(&nilaiModel.Nilai{}).Where("siswa_id = ?", st.ID).Count(&gradeCount)
		if gradeCount < 24 {
			log.Fatalf("Verification failed: student %s (ID %d) has only %d grades (expected at least 24)\n", st.Nama, st.ID, gradeCount)
		}
	}
	log.Println("✓ Check 5 passed: Every student has grades.")

	// 6. User count matches exactly 151 and Administrator Count = 1
	var userCount int64
	db.Model(&authModel.User{}).Count(&userCount)
	if userCount != 151 { // 1 admin + 25 teachers + 125 students = 151
		log.Fatalf("Verification failed: total user count is %d (expected 151)\n", userCount)
	}
	var adminCount int64
	db.Model(&authModel.User{}).Where("role = ?", "admin").Count(&adminCount)
	if adminCount != 1 {
		log.Fatalf("Verification failed: Administrator count is %d (expected exactly 1)\n", adminCount)
	}
	log.Println("✓ Check 6 passed: Total user counts match (151 users) and Administrator Count = 1.")

	// 7. Schedule count matches exactly 150
	var scheduleCount int64
	db.Model(&jadwalModel.Jadwal{}).Count(&scheduleCount)
	if scheduleCount != 150 {
		log.Fatalf("Verification failed: total schedule count is %d (expected exactly 150)\n", scheduleCount)
	}
	log.Println("✓ Check 7 passed: Exactly 150 teaching schedules generated.")

	// 8. Total grades matches exactly 3000
	var totalGrades int64
	db.Model(&nilaiModel.Nilai{}).Count(&totalGrades)
	if totalGrades != 3000 {
		log.Fatalf("Verification failed: total grades count is %d (expected exactly 3000)\n", totalGrades)
	}
	log.Println("✓ Check 8 passed: Exactly 3000 academic grades seeded.")

	log.Println("All automated verification checks passed successfully!")
}
