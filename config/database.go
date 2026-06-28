package config

import (
	"log"
	"os"

	authModel "backend/modules/auth/model"
	guruModel "backend/modules/guru/model"
	jadwalModel "backend/modules/jadwal/model"
	kelasModel "backend/modules/kelas/model"
	mapelModel "backend/modules/mapel/model"
	nilaiModel "backend/modules/nilai/model"
	siswaModel "backend/modules/siswa/model"
	absensiModel "backend/modules/absensi/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {

	err := godotenv.Overload()

	if err != nil {
		log.Fatal("Gagal load .env")
	}

	dsn := os.Getenv("DATABASE_URL")

	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Gagal koneksi database")
	}

	DB = database

	err = DB.AutoMigrate(
		&authModel.User{},
		&guruModel.Guru{},
		&kelasModel.Kelas{},
		&siswaModel.Siswa{},
		&mapelModel.Mapel{},
		&jadwalModel.Jadwal{},
		&nilaiModel.Nilai{},
		&absensiModel.Absensi{},
	)
	if err != nil {
		log.Fatal("Gagal menjalankan auto migrate")
	}

	log.Println("Database berhasil terkoneksi")
}
