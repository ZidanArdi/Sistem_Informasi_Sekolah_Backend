package service

import (
	"errors"
	"strings"

	"backend/modules/mapel/model"
	"backend/modules/mapel/repository"
)

func GetAllMapel(search string) ([]model.Mapel, error) {
	return repository.GetAllMapel(search)
}

func GetMapelByID(id uint) (model.Mapel, error) {
	return repository.GetMapelByID(id)
}

func CreateMapel(data model.Mapel) (model.Mapel, error) {
	if err := validateMapel(data); err != nil {
		return model.Mapel{}, err
	}
	if repository.CheckKodeMapelExists(data.KodeMapel, 0) {
		return model.Mapel{}, errors.New("kode_mapel sudah terdaftar")
	}
	return repository.CreateMapel(data)
}

func UpdateMapel(id uint, data model.Mapel) (model.Mapel, error) {
	if err := validateMapel(data); err != nil {
		return model.Mapel{}, err
	}
	if repository.CheckKodeMapelExists(data.KodeMapel, id) {
		return model.Mapel{}, errors.New("kode_mapel sudah digunakan mapel lain")
	}
	return repository.UpdateMapel(id, data)
}

func DeleteMapel(id uint) error {
	return repository.DeleteMapel(id)
}

func validateMapel(data model.Mapel) error {
	if strings.TrimSpace(data.KodeMapel) == "" ||
		strings.TrimSpace(data.NamaMapel) == "" ||
		data.Jam <= 0 {
		return errors.New("kode_mapel, nama_mapel, dan jam wajib diisi")
	}

	return nil
}
