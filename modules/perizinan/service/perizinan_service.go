package service

import (
	"errors"
	"strings"

	"backend/modules/perizinan/model"
	"backend/modules/perizinan/repository"
)

func CreatePerizinan(data model.Perizinan, currentRole string, currentUserEmail string) (model.Perizinan, error) {
	currentRole = strings.ToLower(strings.TrimSpace(currentRole))
	if currentRole != "siswa" {
		return model.Perizinan{}, errors.New("hanya siswa yang dapat mengajukan perizinan")
	}

	siswaID, err := repository.GetSiswaIDByEmail(currentUserEmail)
	if err != nil {
		return model.Perizinan{}, errors.New("profil siswa tidak ditemukan")
	}

	data.SiswaID = siswaID
	data.Status = "Pending"
	data.DisetujuiOleh = nil
	data.KeteranganGuru = nil

	if data.TanggalMulai.IsZero() || data.TanggalSelesai.IsZero() {
		return model.Perizinan{}, errors.New("tanggal mulai dan tanggal selesai wajib diisi")
	}

	if data.TanggalSelesai.Before(data.TanggalMulai) {
		return model.Perizinan{}, errors.New("tanggal selesai tidak boleh mendahului tanggal mulai")
	}

	if strings.TrimSpace(data.Tipe) == "" {
		return model.Perizinan{}, errors.New("tipe perizinan wajib diisi")
	}

	if strings.TrimSpace(data.Alasan) == "" {
		return model.Perizinan{}, errors.New("alasan perizinan wajib diisi")
	}

	return repository.CreatePerizinan(data)
}

func GetPerizinan(currentRole string, currentUserEmail string) ([]model.Perizinan, error) {
	currentRole = strings.ToLower(strings.TrimSpace(currentRole))

	switch currentRole {
	case "siswa":
		siswaID, err := repository.GetSiswaIDByEmail(currentUserEmail)
		if err != nil {
			return nil, errors.New("profil siswa tidak ditemukan")
		}
		return repository.GetPerizinanBySiswaID(siswaID)
	case "guru":
		// Guru sees all pending requests
		return repository.GetPendingPerizinan()
	case "admin":
		// Admin monitors all requests
		return repository.GetAllPerizinan()
	default:
		return nil, errors.New("role tidak dikenali")
	}
}

func ApproveOrRejectPerizinan(id uint, status string, keteranganGuru string, currentRole string, currentUserEmail string) (model.Perizinan, error) {
	currentRole = strings.ToLower(strings.TrimSpace(currentRole))
	if currentRole != "guru" {
		return model.Perizinan{}, errors.New("hanya guru yang dapat menyetujui atau menolak perizinan")
	}

	status = strings.TrimSpace(status)
	if status != "Disetujui" && status != "Ditolak" {
		return model.Perizinan{}, errors.New("status persetujuan harus Disetujui atau Ditolak")
	}

	guruID, err := repository.GetGuruIDByEmail(currentUserEmail)
	if err != nil {
		return model.Perizinan{}, errors.New("profil guru tidak ditemukan")
	}

	var data model.Perizinan
	data.Status = status
	if strings.TrimSpace(keteranganGuru) != "" {
		data.KeteranganGuru = &keteranganGuru
	}
	data.DisetujuiOleh = &guruID

	return repository.UpdatePerizinan(id, data)
}
