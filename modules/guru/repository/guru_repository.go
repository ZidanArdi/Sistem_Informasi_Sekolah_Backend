package repository

import (
	"fmt"

	"backend/config"
	authModel "backend/modules/auth/model"
	"backend/modules/guru/model"

	"gorm.io/gorm"
)

func BeginTransaction() *gorm.DB {
	return config.DB.Begin()
}

func GetAllGuru(search string) ([]model.Guru, error) {
	var guru []model.Guru
	query := config.DB.Preload("User")
	if search != "" {
		query = query.Where("nama ILIKE ? OR nip ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	result := query.Find(&guru)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(guru) > 0 {
		guruIDs := make([]uint, len(guru))
		for i := range guru {
			guruIDs[i] = guru[i].ID
		}

		var guruMapels []model.GuruMapel
		if err := config.DB.Where("guru_id IN ?", guruIDs).Find(&guruMapels).Error; err == nil {
			mapelMap := make(map[uint][]uint)
			for _, gm := range guruMapels {
				mapelMap[gm.GuruID] = append(mapelMap[gm.GuruID], gm.MapelID)
			}
			for i := range guru {
				guru[i].MapelIDs = mapelMap[guru[i].ID]
			}
		}
	}

	return guru, nil
}

func GetGuruByID(id uint) (model.Guru, error) {
	var guru model.Guru
	result := config.DB.Preload("User").First(&guru, id)
	if result.Error != nil {
		return guru, result.Error
	}

	var mapelIDs []uint
	config.DB.Model(&model.GuruMapel{}).Where("guru_id = ?", guru.ID).Pluck("mapel_id", &mapelIDs)
	guru.MapelIDs = mapelIDs
	return guru, nil
}

func UpdateGuru(id uint, data model.Guru) (model.Guru, error) {
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var guru model.Guru
	if err := tx.First(&guru, id).Error; err != nil {
		tx.Rollback()
		return guru, err
	}

	guru.Nama = data.Nama
	guru.Gelar = data.Gelar
	guru.JenisKelamin = data.JenisKelamin
	guru.NoHP = data.NoHP
	guru.Provinsi = data.Provinsi
	guru.Kabupaten = data.Kabupaten
	guru.Kecamatan = data.Kecamatan
	guru.Desa = data.Desa
	guru.AlamatDetail = data.AlamatDetail

	if err := tx.Save(&guru).Error; err != nil {
		tx.Rollback()
		return guru, err
	}

	if err := tx.Where("guru_id = ?", id).Delete(&model.GuruMapel{}).Error; err != nil {
		tx.Rollback()
		return guru, err
	}

	if len(data.MapelIDs) > 0 {
		for _, mapelID := range data.MapelIDs {
			gm := model.GuruMapel{GuruID: id, MapelID: mapelID}
			if err := tx.Create(&gm).Error; err != nil {
				tx.Rollback()
				return guru, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return guru, err
	}

	config.DB.Preload("User").First(&guru, guru.ID)
	var mapelIDs []uint
	config.DB.Model(&model.GuruMapel{}).Where("guru_id = ?", guru.ID).Pluck("mapel_id", &mapelIDs)
	guru.MapelIDs = mapelIDs
	return guru, nil
}

func DeleteGuru(id uint) error {
	var guru model.Guru
	if err := config.DB.First(&guru, id).Error; err != nil {
		return err
	}
	return config.DB.Delete(&guru).Error
}

func CheckNIPExists(nip string, excludeID uint) bool {
	var count int64
	query := config.DB.Model(&model.Guru{}).Where("nip = ?", nip)
	if excludeID != 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count > 0
}

func IsNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}

func GenerateNIPWithTx(tx *gorm.DB) (string, error) {
	prefix := "GR"
	var latestGuru model.Guru
	err := tx.Set("gorm:query_option", "FOR UPDATE").
		Unscoped().
		Where("nip LIKE ?", prefix+"%").
		Order("nip desc").
		First(&latestGuru).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return prefix + "0001", nil
		}
		return "", fmt.Errorf("ERR_GENERATE_NIG: %v", err)
	}

	var seq int
	_, err = fmt.Sscanf(latestGuru.NIP, prefix+"%04d", &seq)
	if err != nil {
		return prefix + "0001", nil
	}

	return fmt.Sprintf("%s%04d", prefix, seq+1), nil
}

func CountUserByEmail(email string) (int64, error) {
	var count int64
	err := config.DB.Model(&authModel.User{}).Where("email = ?", email).Count(&count).Error
	return count, err
}

func CountUserByEmailWithTx(tx *gorm.DB, email string) (int64, error) {
	var count int64
	err := tx.Model(&authModel.User{}).Where("email = ?", email).Count(&count).Error
	return count, err
}

func CountGuruByNIPWithTx(tx *gorm.DB, nip string) (int64, error) {
	var count int64
	err := tx.Model(&model.Guru{}).Where("nip = ?", nip).Count(&count).Error
	return count, err
}

func CreateUserWithTx(tx *gorm.DB, user *authModel.User) error {
	return tx.Create(user).Error
}

func CreateGuruWithTx(tx *gorm.DB, guru *model.Guru) error {
	return tx.Create(guru).Error
}

func CreateGuruMapelWithTx(tx *gorm.DB, gm *model.GuruMapel) error {
	return tx.Create(gm).Error
}

func LoadGuruRelations(guru *model.Guru) {
	config.DB.Preload("User").First(guru, guru.ID)
	var mapelIDs []uint
	config.DB.Model(&model.GuruMapel{}).Where("guru_id = ?", guru.ID).Pluck("mapel_id", &mapelIDs)
	guru.MapelIDs = mapelIDs
}
