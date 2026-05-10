package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"backend/modules/siswa/model"
	
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

	DB.AutoMigrate(&model.Siswa{})

	log.Println("Database berhasil terkoneksi")
}