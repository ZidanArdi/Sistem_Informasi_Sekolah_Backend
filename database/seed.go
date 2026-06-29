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

	// 1. Seed Guru
	var guru guruModel.Guru
	err = db.Where("email = ?", "guru@sekolah.com").First(&guru).Error
	if err != nil {
		guru = guruModel.Guru{
			NIP:          "198503102010121001",
			Nama:         "Budi Utomo, S.Pd.",
			JenisKelamin: "Laki-laki",
			Email:        "guru@sekolah.com",
			NoHP:         "081234567890",
			Alamat:       "Jl. Cendrawasih No. 10, Semarang",
		}
		if err := db.Create(&guru).Error; err != nil {
			log.Fatal("Gagal seeding Guru:", err)
		}
		log.Println("Guru berhasil diseed.")
	} else {
		log.Println("Guru sudah ada di database.")
	}

	// 2. Seed Kelas
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

	// 3. Seed Siswa
	var siswa siswaModel.Siswa
	err = db.Where("email = ?", "siswa@sekolah.com").First(&siswa).Error
	if err != nil {
		siswa = siswaModel.Siswa{
			NIS:          "10122045",
			Nama:         "Rian Hidayat",
			JenisKelamin: "Laki-laki",
			TempatLahir:  "Semarang",
			TanggalLahir: "2008-05-12",
			Alamat:       "Jl. Cempaka Raya No. 45, Semarang",
			Email:        "siswa@sekolah.com",
			KelasID:      kelas.ID,
		}
		if err := db.Create(&siswa).Error; err != nil {
			log.Fatal("Gagal seeding Siswa:", err)
		}
		log.Println("Siswa berhasil diseed.")
	} else {
		log.Println("Siswa sudah ada di database.")
	}

	// 4. Seed Users
	users := []authModel.User{
		{
			Nama:     "Admin Sistem",
			Email:    "admin@sekolah.com",
			Password: string(hashedPassword),
			Role:     "admin",
		},
		{
			Nama:     "Siti Aminah (Staff TU)",
			Email:    "staff@sekolah.com",
			Password: string(hashedPassword),
			Role:     "staff",
		},
		{
			Nama:     "Budi Utomo, S.Pd.",
			Email:    "guru@sekolah.com",
			Password: string(hashedPassword),
			Role:     "guru",
		},
		{
			Nama:     "Rian Hidayat",
			Email:    "siswa@sekolah.com",
			Password: string(hashedPassword),
			Role:     "siswa",
		},
	}

	for _, user := range users {
		var existing authModel.User
		err = db.Where("email = ?", user.Email).First(&existing).Error
		if err != nil {
			if err := db.Create(&user).Error; err != nil {
				log.Fatal("Gagal seeding User:", user.Email, err)
			}
			log.Printf("User %s (%s) berhasil diseed.\n", user.Nama, user.Role)
		} else {
			log.Printf("User %s sudah ada di database.\n", user.Email)
		}
	}

	log.Println("Seeding database selesai dengan sukses!")
}
