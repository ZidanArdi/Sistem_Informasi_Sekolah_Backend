package service

import (
	"errors"
	"strings"

	"backend/modules/guru/model"
	"backend/modules/guru/repository"
)

func GetAllGuru(search string) ([]model.Guru, error) {
	return repository.GetAllGuru(search)
}

func GetGuruByID(id uint) (model.Guru, error) {
	return repository.GetGuruByID(id)
}

func CreateGuru(data model.Guru) (model.Guru, error) {
	if err := validateGuru(data); err != nil {
		return model.Guru{}, err
	}
	if repository.CheckNIPExists(data.NIP, 0) {
		return model.Guru{}, errors.New("NIP sudah terdaftar")
	}
	return repository.CreateGuru(data)
}

func UpdateGuru(id uint, data model.Guru) (model.Guru, error) {
	if err := validateGuru(data); err != nil {
		return model.Guru{}, err
	}
	if repository.CheckNIPExists(data.NIP, id) {
		return model.Guru{}, errors.New("NIP sudah digunakan guru lain")
	}
	return repository.UpdateGuru(id, data)
}

func DeleteGuru(id uint) error {
	return repository.DeleteGuru(id)
}

func validateGuru(data model.Guru) error {
	if strings.TrimSpace(data.NIP) == "" ||
		strings.TrimSpace(data.Nama) == "" ||
		strings.TrimSpace(data.JenisKelamin) == "" {
		return errors.New("nip, nama, dan jenis_kelamin wajib diisi")
	}

	return nil
}
