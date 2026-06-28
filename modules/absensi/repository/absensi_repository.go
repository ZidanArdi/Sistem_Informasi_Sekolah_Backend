package repository

import (
	"backend/config"
	"backend/modules/absensi/model"
	"strconv"

	"gorm.io/gorm"
)

func GetAllAbsensi(kelasID string, tanggal string, statusPersetujuan string, siswaID uint) ([]model.Absensi, error) {
	var absensiList []model.Absensi

	query := config.DB.Preload("Siswa").Preload("Siswa.Kelas")

	if siswaID != 0 {
		query = query.Where("siswa_id = ?", siswaID)
	}

	if tanggal != "" {
		query = query.Where("tanggal = ?", tanggal)
	}

	if statusPersetujuan != "" {
		query = query.Where("status_persetujuan = ?", statusPersetujuan)
	}

	if kelasID != "" {
		if parsedKelasID, err := strconv.Atoi(kelasID); err == nil {
			// Join with siswa table to filter by kelas_id
			query = query.Joins("JOIN siswa ON siswa.id = absensi.siswa_id").Where("siswa.kelas_id = ?", parsedKelasID)
		}
	}

	result := query.Order("tanggal desc, id desc").Find(&absensiList)
	return absensiList, result.Error
}

func GetAbsensiByID(id uint) (model.Absensi, error) {
	var absensi model.Absensi
	result := config.DB.Preload("Siswa").Preload("Siswa.Kelas").First(&absensi, id)
	return absensi, result.Error
}

func CreateAbsensi(data model.Absensi) (model.Absensi, error) {
	result := config.DB.Create(&data)
	return data, result.Error
}

func UpdateAbsensi(id uint, data model.Absensi) (model.Absensi, error) {
	var absensi model.Absensi
	err := config.DB.First(&absensi, id).Error
	if err != nil {
		return absensi, err
	}

	absensi.Status = data.Status
	absensi.Keterangan = data.Keterangan
	absensi.StatusPersetujuan = data.StatusPersetujuan

	if err := config.DB.Save(&absensi).Error; err != nil {
		return absensi, err
	}

	config.DB.Preload("Siswa").Preload("Siswa.Kelas").First(&absensi, absensi.ID)
	return absensi, nil
}

func DeleteAbsensi(id uint) error {
	return config.DB.Delete(&model.Absensi{}, id).Error
}

func BulkSaveAbsensi(records []model.Absensi) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		for _, record := range records {
			var existing model.Absensi
			err := tx.Where("siswa_id = ? AND tanggal = ?", record.SiswaID, record.Tanggal).First(&existing).Error
			if err == nil {
				// Record exists, update it
				existing.Status = record.Status
				existing.Keterangan = record.Keterangan
				existing.StatusPersetujuan = record.StatusPersetujuan
				if err := tx.Save(&existing).Error; err != nil {
					return err
				}
			} else {
				// Record does not exist, create new
				if err := tx.Create(&record).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func GetSiswaIDByEmail(email string) (uint, error) {
	var siswa struct {
		ID uint
	}
	err := config.DB.Table("siswa").Select("id").Where("email = ? AND deleted_at IS NULL", email).First(&siswa).Error
	if err != nil {
		return 0, err
	}
	return siswa.ID, nil
}

func IsNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}
