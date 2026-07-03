package main

import (
	"log"

	"backend/config"
	authModel "backend/modules/auth/model"
	guruModel "backend/modules/guru/model"
	kelasModel "backend/modules/kelas/model"
	siswaModel "backend/modules/siswa/model"

	"golang.org/x/crypto/bcrypt"
)

func strPtr(s string) *string {
	return &s
}

func main() {
	// Koneksi ke Database
	config.ConnectDB()

	db := config.DB
	log.Println("Memulai proses seeding...")

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Gagal melakukan hashing password:", err)
	}

	// 1. Seed Users (Terlebih dahulu agar bisa direlasikan)
	users := []authModel.User{
		{
			Nama:         "Admin Sistem",
			Email:        strPtr("admin@sekolah.com"),
			Password:     string(hashedPassword),
			Role:         "admin",
			IsFirstLogin: false,
		},
		{
			Nama:         "Siti Aminah (Admin TU)",
			Email:        strPtr("staff@sekolah.com"),
			Password:     string(hashedPassword),
			Role:         "admin",
			IsFirstLogin: false,
		},
		{
			Nama:         "Budi Utomo, S.Pd.",
			Email:        strPtr("guru@sekolah.com"),
			Password:     string(hashedPassword),
			Role:         "guru",
			IsFirstLogin: false,
		},
		{
			Nama:         "Rian Hidayat",
			Email:        strPtr("siswa@sekolah.com"),
			Password:     string(hashedPassword),
			Role:         "siswa",
			IsFirstLogin: false,
		},
	}

	for _, user := range users {
		var existing authModel.User
		err = db.Where("email = ?", *user.Email).First(&existing).Error
		if err != nil {
			if err := db.Create(&user).Error; err != nil {
				log.Fatal("Gagal seeding User:", *user.Email, err)
			}
			log.Printf("User %s (%s) berhasil diseed.\n", user.Nama, user.Role)
		} else {
			existing.Role = user.Role
			existing.Nama = user.Nama
			existing.IsFirstLogin = user.IsFirstLogin
			db.Save(&existing)
			log.Printf("User %s berhasil diupdate (role: %s).\n", *user.Email, user.Role)
		}
	}

	// Dapatkan User Guru Budi Utomo
	var guruUser authModel.User
	err = db.Where("email = ?", "guru@sekolah.com").First(&guruUser).Error
	if err != nil {
		log.Fatal("User guru@sekolah.com tidak ditemukan untuk relasi Guru")
	}

	// 2. Seed Guru
	var guru guruModel.Guru
	err = db.Where("nip = ?", "198503102010121001").First(&guru).Error
	if err != nil {
		guru = guruModel.Guru{
			UserID:       guruUser.ID,
			NIP:          "198503102010121001",
			Nama:         "Budi Utomo",
			Gelar:        "S.Pd.",
			JenisKelamin: "Laki-laki",
			NoHP:         "081234567890",
			Alamat:       "Jl. Cendrawasih No. 10, Semarang",
		}
		if err := db.Create(&guru).Error; err != nil {
			log.Fatal("Gagal seeding Guru:", err)
		}
		log.Println("Guru Budi Utomo berhasil diseed.")
	} else {
		if guru.UserID == 0 {
			guru.UserID = guruUser.ID
			db.Save(&guru)
			log.Println("Hubungan user_id guru Budi Utomo berhasil diperbarui.")
		}
		log.Println("Guru sudah ada di database.")
	}

	// 3. Seed Kelas
	var kelas kelasModel.Kelas
	err = db.Where("nama_kelas = ?", "XI-MIPA-1").First(&kelas).Error
	if err != nil {
		kelas = kelasModel.Kelas{
			NamaKelas:   "XI-MIPA-1",
			Tingkat:     "11",
			WaliKelasID: &guru.ID,
		}
		if err := db.Create(&kelas).Error; err != nil {
			log.Fatal("Gagal seeding Kelas:", err)
		}
		log.Println("Kelas berhasil diseed.")
	} else {
		log.Println("Kelas sudah ada di database.")
	}

	// Dapatkan User Siswa Rian Hidayat
	var siswaUser authModel.User
	err = db.Where("email = ?", "siswa@sekolah.com").First(&siswaUser).Error
	if err != nil {
		log.Fatal("User siswa@sekolah.com tidak ditemukan untuk relasi Siswa")
	}

	// 4. Seed Siswa
	var siswa siswaModel.Siswa
	err = db.Where("email = ?", "siswa@sekolah.com").First(&siswa).Error
	if err != nil {
		siswa = siswaModel.Siswa{
			UserID:       siswaUser.ID,
			NIS:          "10122045",
			Nama:         "Rian Hidayat",
			JenisKelamin: "Laki-laki",
			TempatLahir:  "Semarang",
			TanggalLahir: "2008-05-12",
			Alamat:       "Jl. Cempaka Raya No. 45, Semarang",
			NoHP:         "082134567890",
			Email:        "siswa@sekolah.com",
			KelasID:      kelas.ID,
		}
		if err := db.Create(&siswa).Error; err != nil {
			log.Fatal("Gagal seeding Siswa:", err)
		}
		log.Println("Siswa berhasil diseed.")
	} else {
		if siswa.UserID == 0 {
			siswa.UserID = siswaUser.ID
			db.Save(&siswa)
			log.Println("Hubungan user_id siswa Rian Hidayat berhasil diperbarui.")
		}
		log.Println("Siswa sudah ada di database.")
	}

	log.Println("Seeding database selesai dengan sukses!")
}
