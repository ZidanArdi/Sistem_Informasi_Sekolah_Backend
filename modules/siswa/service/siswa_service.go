package service

import (
	"errors"
	"strings"

	"backend/modules/siswa/model"
	"backend/modules/siswa/repository"
)

func GetAllSiswa(search string, kelasID string) ([]model.Siswa, error) {
	return repository.GetAllSiswa(search, kelasID)
}

func GetSiswaByID(id uint) (model.Siswa, error) {
	return repository.GetSiswaByID(id)
}

func CreateSiswa(data model.Siswa) (model.Siswa, error) {
	if err := validateSiswa(data); err != nil {
		return model.Siswa{}, err
	}
	if repository.CheckNISExists(data.NIS, 0) {
		return model.Siswa{}, errors.New("NIS sudah terdaftar")
	}
	return repository.CreateSiswa(data)
}

func UpdateSiswa(id uint, data model.Siswa) (model.Siswa, error) {
	if err := validateSiswa(data); err != nil {
		return model.Siswa{}, err
	}
	if repository.CheckNISExists(data.NIS, id) {
		return model.Siswa{}, errors.New("NIS sudah digunakan siswa lain")
	}
	return repository.UpdateSiswa(id, data)
}

func DeleteSiswa(id uint) error {
	return repository.DeleteSiswa(id)
}

func validateSiswa(data model.Siswa) error {
	if strings.TrimSpace(data.NIS) == "" ||
		strings.TrimSpace(data.Nama) == "" ||
		strings.TrimSpace(data.JenisKelamin) == "" ||
		strings.TrimSpace(data.TempatLahir) == "" ||
		strings.TrimSpace(data.TanggalLahir) == "" ||
		strings.TrimSpace(data.Alamat) == "" ||
		data.KelasID == 0 {
		return errors.New("nis, nama, jenis_kelamin, tempat_lahir, tanggal_lahir, alamat, dan kelas_id wajib diisi")
	}

	return nil
}
