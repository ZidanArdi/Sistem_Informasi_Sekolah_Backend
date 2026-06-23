package service

import (
	"errors"
	"strings"

	"backend/modules/nilai/model"
	"backend/modules/nilai/repository"
)

func GetAllNilai(siswaID string, mapelID string, semester string, jenisNilai string) ([]model.Nilai, error) {
	return repository.GetAllNilai(siswaID, mapelID, semester, jenisNilai)
}

func GetNilaiByID(id uint) (model.Nilai, error) {
	return repository.GetNilaiByID(id)
}

func CreateNilai(data model.Nilai) (model.Nilai, error) {
	if err := validateNilai(data); err != nil {
		return model.Nilai{}, err
	}
	return repository.CreateNilai(data)
}

func UpdateNilai(id uint, data model.Nilai) (model.Nilai, error) {
	if err := validateNilai(data); err != nil {
		return model.Nilai{}, err
	}
	return repository.UpdateNilai(id, data)
}

func DeleteNilai(id uint) error {
	return repository.DeleteNilai(id)
}

func validateNilai(data model.Nilai) error {
	jenisNilai := strings.TrimSpace(strings.ToLower(data.JenisNilai))
	if data.SiswaID == 0 ||
		data.MapelID == 0 ||
		strings.TrimSpace(data.Semester) == "" ||
		jenisNilai == "" {
		return errors.New("siswa_id, mapel_id, semester, jenis_nilai, dan nilai wajib diisi")
	}

	if jenisNilai != "tugas" && jenisNilai != "uts" && jenisNilai != "uas" {
		return errors.New("jenis_nilai harus tugas, uts, atau uas")
	}

	if data.Nilai < 0 || data.Nilai > 100 {
		return errors.New("nilai harus di antara 0 dan 100")
	}

	data.JenisNilai = jenisNilai
	return nil
}
