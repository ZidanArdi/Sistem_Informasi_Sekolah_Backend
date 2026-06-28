package service

import (
	"errors"
	"strings"

	"backend/modules/absensi/model"
	"backend/modules/absensi/repository"
)

func GetAllAbsensi(kelasID string, tanggal string, statusPersetujuan string, specificSiswaID uint, currentRole string, currentUserEmail string) ([]model.Absensi, error) {
	var siswaID uint = specificSiswaID

	if strings.ToLower(currentRole) == "siswa" {
		id, err := repository.GetSiswaIDByEmail(currentUserEmail)
		if err != nil {
			return nil, errors.New("siswa dengan email ini tidak ditemukan")
		}
		siswaID = id
	}

	return repository.GetAllAbsensi(kelasID, tanggal, statusPersetujuan, siswaID)
}

func GetAbsensiByID(id uint) (model.Absensi, error) {
	return repository.GetAbsensiByID(id)
}

func CreateAbsensi(data model.Absensi, currentRole string, currentUserEmail string) (model.Absensi, error) {
	if data.Tanggal == "" {
		return model.Absensi{}, errors.New("tanggal wajib diisi")
	}
	if data.Status == "" {
		return model.Absensi{}, errors.New("status absensi wajib diisi")
	}

	if strings.ToLower(currentRole) == "siswa" {
		id, err := repository.GetSiswaIDByEmail(currentUserEmail)
		if err != nil {
			return model.Absensi{}, errors.New("siswa dengan email ini tidak ditemukan")
		}
		data.SiswaID = id
		data.StatusPersetujuan = "Pending"
	} else {
		if data.SiswaID == 0 {
			return model.Absensi{}, errors.New("siswa_id wajib diisi")
		}
		if data.StatusPersetujuan == "" {
			data.StatusPersetujuan = "Disetujui"
		}
	}

	return repository.CreateAbsensi(data)
}

func BulkSaveAbsensi(records []model.Absensi, tanggal string) error {
	if tanggal == "" {
		return errors.New("tanggal wajib diisi")
	}

	for i := range records {
		records[i].Tanggal = tanggal
		if records[i].StatusPersetujuan == "" {
			records[i].StatusPersetujuan = "Disetujui"
		}
		if records[i].Status == "" {
			records[i].Status = "Hadir"
		}
	}

	return repository.BulkSaveAbsensi(records)
}

func UpdateAbsensi(id uint, data model.Absensi) (model.Absensi, error) {
	return repository.UpdateAbsensi(id, data)
}

func DeleteAbsensi(id uint) error {
	return repository.DeleteAbsensi(id)
}

func ApproveOrRejectAbsensi(id uint, statusPersetujuan string) (model.Absensi, error) {
	statusPersetujuan = strings.TrimSpace(statusPersetujuan)
	if statusPersetujuan != "Disetujui" && statusPersetujuan != "Ditolak" {
		return model.Absensi{}, errors.New("status persetujuan harus Disetujui atau Ditolak")
	}

	existing, err := repository.GetAbsensiByID(id)
	if err != nil {
		return model.Absensi{}, err
	}

	existing.StatusPersetujuan = statusPersetujuan
	return repository.UpdateAbsensi(id, existing)
}
