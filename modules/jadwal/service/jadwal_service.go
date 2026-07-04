package service

import (
	"errors"
	"strconv"
	"strings"

	"backend/config"
	"backend/modules/jadwal/model"
	"backend/modules/jadwal/repository"
)

func GetAllJadwal(kelasID string, mapelID string, guruID string, hari string, tahunAjaran string, semester string) ([]model.Jadwal, error) {
	return repository.GetAllJadwal(kelasID, mapelID, guruID, hari, tahunAjaran, semester)
}

func GetJadwalByID(id uint) (model.Jadwal, error) {
	return repository.GetJadwalByID(id)
}

func CreateJadwal(data model.Jadwal) (model.Jadwal, error) {
	if err := validateJadwal(data); err != nil {
		return model.Jadwal{}, err
	}
	if err := CheckScheduleConflicts(data); err != nil {
		return model.Jadwal{}, err
	}
	return repository.CreateJadwal(data)
}

func UpdateJadwal(id uint, data model.Jadwal) (model.Jadwal, error) {
	data.ID = id
	if err := validateJadwal(data); err != nil {
		return model.Jadwal{}, err
	}
	if err := CheckScheduleConflicts(data); err != nil {
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
		strings.TrimSpace(data.JamSelesai) == "" ||
		strings.TrimSpace(data.TahunAjaran) == "" ||
		strings.TrimSpace(data.Semester) == "" {
		return errors.New("kelas_id, mapel_id, guru_id, hari, jam_mulai, jam_selesai, tahun_ajaran, dan semester wajib diisi")
	}

	// 1. Teacher Eligibility Validation (check guru_mapel relation)
	var count int64
	err := config.DB.Table("guru_mapel").
		Where("guru_id = ? AND mapel_id = ?", data.GuruID, data.MapelID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Guru tidak dapat mengajar mata pelajaran ini.")
	}

	return nil
}

func timeToMinutes(t string) (int, error) {
	t = strings.TrimSpace(t)
	parts := strings.Split(t, ":")
	if len(parts) < 2 {
		return 0, errors.New("format waktu tidak valid")
	}
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	return hours*60 + minutes, nil
}

func CheckScheduleConflicts(data model.Jadwal) error {
	startMin, err := timeToMinutes(data.JamMulai)
	if err != nil {
		return errors.New("format jam mulai tidak valid")
	}
	endMin, err := timeToMinutes(data.JamSelesai)
	if err != nil {
		return errors.New("format jam selesai tidak valid")
	}
	if startMin >= endMin {
		return errors.New("jam mulai harus sebelum jam selesai")
	}

	var conflicts []model.Jadwal
	query := config.DB.Model(&model.Jadwal{}).
		Where("hari = ? AND tahun_ajaran = ? AND semester = ?", data.Hari, data.TahunAjaran, data.Semester)

	if data.ID != 0 {
		query = query.Where("id != ?", data.ID)
	}

	err = query.Where("guru_id = ? OR kelas_id = ?", data.GuruID, data.KelasID).Find(&conflicts).Error
	if err != nil {
		return err
	}

	for _, c := range conflicts {
		cStart, err := timeToMinutes(c.JamMulai)
		if err != nil {
			continue
		}
		cEnd, err := timeToMinutes(c.JamSelesai)
		if err != nil {
			continue
		}

		// Overlap condition: startA < endB && startB < endA
		if startMin < cEnd && cStart < endMin {
			if c.GuruID == data.GuruID {
				return errors.New("konflik jadwal: Guru tersebut sudah memiliki jadwal mengajar di kelas lain pada waktu yang sama")
			}
			if c.KelasID == data.KelasID {
				return errors.New("konflik jadwal: Kelas tersebut sudah memiliki mata pelajaran lain pada waktu yang sama")
			}
		}
	}

	return nil
}
