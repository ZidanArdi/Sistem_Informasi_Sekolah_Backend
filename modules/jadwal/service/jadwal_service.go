package service

import (
	"errors"
	"strings"

	"backend/modules/jadwal/model"
	"backend/modules/jadwal/repository"
)

func GetAllJadwal(kelasID string, mapelID string, guruID string, hari string) ([]model.Jadwal, error) {
	return repository.GetAllJadwal(kelasID, mapelID, guruID, hari)
}

func GetJadwalByID(id uint) (model.Jadwal, error) {
	return repository.GetJadwalByID(id)
}

func CreateJadwal(data model.Jadwal) (model.Jadwal, error) {
	if err := validateJadwal(data); err != nil {
		return model.Jadwal{}, err
	}
	return repository.CreateJadwal(data)
}

func UpdateJadwal(id uint, data model.Jadwal) (model.Jadwal, error) {
	if err := validateJadwal(data); err != nil {
		return model.Jadwal{}, err
	}
	return repository.UpdateJadwal(id, data)
}

func DeleteJadwal(id uint) error {
	return repository.DeleteJadwal(id)
}

func validateJadwal(data model.Jadwal) error {
	if data.KelasID == 0 ||
		data.MapelID == 0 ||
		data.GuruID == 0 ||
		strings.TrimSpace(data.Hari) == "" ||
		strings.TrimSpace(data.JamMulai) == "" ||
		strings.TrimSpace(data.JamSelesai) == "" {
		return errors.New("kelas_id, mapel_id, guru_id, hari, jam_mulai, dan jam_selesai wajib diisi")
	}

	return nil
}
