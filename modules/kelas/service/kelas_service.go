package service

import (
	"errors"
	"strings"

	"backend/modules/kelas/model"
	"backend/modules/kelas/repository"
)

func GetAllKelas(search string, tingkat string) ([]model.Kelas, error) {
	return repository.GetAllKelas(search, tingkat)
}

func GetKelasByID(id uint) (model.Kelas, error) {
	return repository.GetKelasByID(id)
}

func CreateKelas(data model.Kelas) (model.Kelas, error) {
	if err := validateKelas(data); err != nil {
		return model.Kelas{}, err
	}
	return repository.CreateKelas(data)
}

func UpdateKelas(id uint, data model.Kelas) (model.Kelas, error) {
	data.ID = id
	if err := validateKelas(data); err != nil {
		return model.Kelas{}, err
	}
	return repository.UpdateKelas(id, data)
}

func DeleteKelas(id uint) error {
	return repository.DeleteKelas(id)
}

func validateKelas(data model.Kelas) error {
	if strings.TrimSpace(data.NamaKelas) == "" ||
		strings.TrimSpace(data.Tingkat) == "" ||
		strings.TrimSpace(data.Jurusan) == "" ||
		data.Kapasitas <= 0 {
		return errors.New("nama_kelas, tingkat, jurusan, dan kapasitas wajib diisi dengan benar")
	}

	// Validate Jurusan enum values
	allowedJurusan := map[string]bool{
		"RPL": true,
		"TKJ": true,
		"AKL": true,
		"DKV": true,
	}
	if !allowedJurusan[data.Jurusan] {
		return errors.New("jurusan tidak valid. Harus salah satu dari: RPL, TKJ, AKL, DKV")
	}

	// Validate WaliKelas assignment (One guru can become wali kelas for one class only)
	if data.WaliKelasID != nil && *data.WaliKelasID != 0 {
		if repository.CheckWaliKelasExists(*data.WaliKelasID, data.ID) {
			return errors.New("guru ini sudah menjadi wali kelas di kelas lain")
		}
	}

	return nil
}
