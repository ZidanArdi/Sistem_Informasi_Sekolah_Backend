package main

import (
	"log"
	"backend/config"
	jadwalModel "backend/modules/jadwal/model"
)

func main() {
	config.ConnectDB()
	var schedules []jadwalModel.Jadwal
	config.DB.Preload("Mapel").Preload("Guru").Preload("Kelas").Find(&schedules)

	log.Printf("Total schedules in DB: %d\n", len(schedules))
	for _, s := range schedules {
		if s.Kelas.NamaKelas == "XI-DKV-1" {
			log.Printf("Class: %s | Day: %s | Slot: %s - %s | Mapel: %s (ID %d) | Guru: %s (ID %d)\n", 
				s.Kelas.NamaKelas, s.Hari, s.JamMulai, s.JamSelesai, s.Mapel.NamaMapel, s.Mapel.ID, s.Guru.Nama, s.Guru.ID)
		}
	}
}
