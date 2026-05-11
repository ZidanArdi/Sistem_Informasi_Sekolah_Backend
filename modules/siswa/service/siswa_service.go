package service

import (
	"backend/modules/siswa/model"
	"backend/modules/siswa/repository"
)

func GetAllSiswa(search string) ([]model.Siswa, error) {
	return repository.GetAllSiswa(search)
}

func GetSiswaByID(id string) (model.Siswa, error) {
	return repository.GetSiswaByID(id)
}

func CreateSiswa(data model.Siswa) (model.Siswa, error) {
	return repository.CreateSiswa(data)
}

func UpdateSiswa(id string, data model.Siswa) (model.Siswa, error) {
	return repository.UpdateSiswa(id, data)
}

func DeleteSiswa(id string) error {
	return repository.DeleteSiswa(id)
}

func CheckEmailExists(email string, excludeNIS string) bool {
	return repository.CheckEmailExists(email, excludeNIS)
}