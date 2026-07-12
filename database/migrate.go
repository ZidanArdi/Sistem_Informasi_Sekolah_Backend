package main

import (
	"log"

	"backend/config"
)

func main() {
	// connect database and run default AutoMigrate first
	config.ConnectDB()
	db := config.DB

	log.Println("Running custom database migrations...")

	// 1. Alter jadwals columns to match time datatype
	log.Println("Converting jadwals time columns to time datatype...")
	if err := db.Exec("ALTER TABLE jadwals ALTER COLUMN jam_mulai TYPE time USING jam_mulai::time").Error; err != nil {
		log.Println("Note/Warning: Altering jam_mulai error:", err)
	}
	if err := db.Exec("ALTER TABLE jadwals ALTER COLUMN jam_selesai TYPE time USING jam_selesai::time").Error; err != nil {
		log.Println("Note/Warning: Altering jam_selesai error:", err)
	}

	// 2. Drop constraints in nilais table
	log.Println("Dropping unused NOT NULL constraints in nilais table...")
	if err := db.Exec("ALTER TABLE nilais ALTER COLUMN jenis_nilai DROP NOT NULL").Error; err != nil {
		log.Println("Note/Warning: Dropping jenis_nilai NOT NULL constraint error (may already be dropped):", err)
	}
	if err := db.Exec("ALTER TABLE nilais ALTER COLUMN nilai DROP NOT NULL").Error; err != nil {
		log.Println("Note/Warning: Dropping nilai NOT NULL constraint error (may already be dropped):", err)
	}

	log.Println("Database migrations completed successfully!")
}
