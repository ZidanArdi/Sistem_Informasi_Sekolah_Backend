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
	if err := validateKelas(data); err != nil {
		return model.Kelas{}, err
	}
	return repository.UpdateKelas(id, data)
}

func DeleteKelas(id uint) error {
	return repository.DeleteKelas(id)
}

func validateKelas(data model.Kelas) error {
	if strings.TrimSpace(data.NamaKelas) == "" || strings.TrimSpace(data.Tingkat) == "" {
		return errors.New("nama_kelas dan tingkat wajib diisi")
	}

	return nil
}
