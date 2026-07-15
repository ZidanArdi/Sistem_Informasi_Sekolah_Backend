package service

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"backend/config"
	authModel "backend/modules/auth/model"
	"backend/modules/guru/model"
	"backend/modules/guru/repository"

	"golang.org/x/crypto/bcrypt"
)

func GetAllGuru(search string) ([]model.Guru, error) {
	return repository.GetAllGuru(search)
}

func GetGuruByID(id uint) (model.Guru, error) {
	return repository.GetGuruByID(id)
}

func GenerateGuruEmail(nama string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
	if err != nil {
		return "", err
	}
	cleanName := reg.ReplaceAllString(nama, "")

	parts := strings.Fields(strings.ToLower(cleanName))

	var nameParts []string
	excludeTitles := map[string]bool{
		"spd": true, "mpd": true, "skom": true, "mkom": true, "dr": true, "prof": true,
	}
	for _, p := range parts {
		if !excludeTitles[p] {
			nameParts = append(nameParts, p)
		}
	}

	baseEmail := strings.Join(nameParts, ".")
	if baseEmail == "" {
		baseEmail = "guru"
	}

	emailDomain := "@sekolah.com"
	email := baseEmail + emailDomain

	count, err := repository.CountUserByEmail(email)
	if err != nil {
		return "", err
	}
	if count == 0 {
		return email, nil
	}

	suffix := 1
	for {
		emailWithSuffix := fmt.Sprintf("%s%02d%s", baseEmail, suffix, emailDomain)
		count, err = repository.CountUserByEmail(emailWithSuffix)
		if err != nil {
			return "", err
		}
		if count == 0 {
			return emailWithSuffix, nil
		}
		suffix++
	}
}

func CreateGuru(data model.Guru) (model.Guru, string, error) {
	tx := repository.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	nip, err := repository.GenerateNIPWithTx(tx)
	if err != nil {
		tx.Rollback()
		return model.Guru{}, "", err
	}
	data.NIP = nip

	var email string
	var seq int
	if len(data.NIP) >= 4 {
		fmt.Sscanf(data.NIP, "GR%d", &seq)
	}
	if seq == 0 {
		seq = 1
	}

	if seq == 1 {
		email = "principal@sekolah.com"
	} else {
		email = fmt.Sprintf("guru%02d@sekolah.com", seq)
	}

	if err := validateGuru(data); err != nil {
		tx.Rollback()
		return model.Guru{}, "", err
	}

	exists, err := repository.CountGuruByNIPWithTx(tx, data.NIP)
	if err != nil {
		tx.Rollback()
		return model.Guru{}, "", fmt.Errorf("ERR_GENERATE_NIG: Gagal memeriksa keunikan NIG: %v", err)
	}
	if exists > 0 {
		tx.Rollback()
		return model.Guru{}, "", errors.New("ERR_IDENTIFIER_DUPLICATE: NIG sudah terdaftar")
	}

	emailCount, err := repository.CountUserByEmailWithTx(tx, email)
	if err != nil {
		tx.Rollback()
		return model.Guru{}, "", fmt.Errorf("ERR_GENERATE_NIG: Gagal memeriksa keunikan email: %v", err)
	}
	if emailCount > 0 {
		tx.Rollback()
		return model.Guru{}, "", errors.New("ERR_EMAIL_DUPLICATE: Email sudah terdaftar")
	}

	plainPassword := "Guru123!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return model.Guru{}, "", err
	}

	user := authModel.User{
		Nama:         data.Nama,
		Email:        &email,
		Password:     string(hashedPassword),
		Role:         "guru",
		IsFirstLogin: true,
	}

	if err := repository.CreateUserWithTx(tx, &user); err != nil {
		tx.Rollback()
		return model.Guru{}, "", errors.New("ERR_EMAIL_DUPLICATE: gagal membuat akun user untuk guru: " + err.Error())
	}

	data.UserID = user.ID
	if err := repository.CreateGuruWithTx(tx, &data); err != nil {
		tx.Rollback()
		return model.Guru{}, "", errors.New("ERR_GENERATE_NIG: gagal menyimpan data guru: " + err.Error())
	}

	if len(data.MapelIDs) > 0 {
		for _, mapelID := range data.MapelIDs {
			gm := model.GuruMapel{GuruID: data.ID, MapelID: mapelID}
			if err := repository.CreateGuruMapelWithTx(tx, &gm); err != nil {
				tx.Rollback()
				return model.Guru{}, "", errors.New("ERR_GENERATE_NIG: gagal menyimpan relasi guru mapel: " + err.Error())
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return model.Guru{}, "", err
	}

	if config.Debug {
		log.Printf("[INFO] Action: Teacher Created | Timestamp: %s | User: %s | Identifier: %s | Email: %s", time.Now().Format(time.RFC3339), data.Nama, data.NIP, email)
	}

	repository.LoadGuruRelations(&data)
	return data, plainPassword, nil
}

func UpdateGuru(id uint, data model.Guru) (model.Guru, error) {
	if err := validateGuru(data); err != nil {
		return model.Guru{}, err
	}
	return repository.UpdateGuru(id, data)
}

func DeleteGuru(id uint) error {
	return repository.DeleteGuru(id)
}

func validateGuru(data model.Guru) error {
	if strings.TrimSpace(data.Nama) == "" ||
		strings.TrimSpace(data.Gelar) == "" ||
		strings.TrimSpace(data.JenisKelamin) == "" ||
		strings.TrimSpace(data.Provinsi) == "" ||
		strings.TrimSpace(data.Kabupaten) == "" ||
		strings.TrimSpace(data.Kecamatan) == "" ||
		strings.TrimSpace(data.Desa) == "" ||
		strings.TrimSpace(data.AlamatDetail) == "" {
		return errors.New("nama, gelar, jenis_kelamin, provinsi, kabupaten, kecamatan, desa, dan alamat_detail wajib diisi")
	}

	return nil
}
