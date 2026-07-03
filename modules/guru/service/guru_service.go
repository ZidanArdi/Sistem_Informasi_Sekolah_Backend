package service

import (
	"errors"
	"fmt"
	"math/rand"
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

	var count int64
	err = config.DB.Model(&authModel.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return "", err
	}
	if count == 0 {
		return email, nil
	}

	suffix := 1
	for {
		emailWithSuffix := fmt.Sprintf("%s%02d%s", baseEmail, suffix, emailDomain)
		err = config.DB.Model(&authModel.User{}).Where("email = ?", emailWithSuffix).Count(&count).Error
		if err != nil {
			return "", err
		}
		if count == 0 {
			return emailWithSuffix, nil
		}
		suffix++
	}
}

func GenerateRandomPassword() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(90000) + 10000
	return fmt.Sprintf("GR-%d", num)
}

func CreateGuru(data model.Guru) (model.Guru, string, error) {
	// Start Database Transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Auto generate NIP with transaction lock
	nip, err := repository.GenerateNIPWithTx(tx)
	if err != nil {
		tx.Rollback()
		return model.Guru{}, "", err
	}
	data.NIP = nip

	// 2. Auto generate Email
	email, err := GenerateGuruEmail(data.Nama)
	if err != nil {
		tx.Rollback()
		return model.Guru{}, "", err
	}

	// Validate fields
	if err := validateGuru(data); err != nil {
		tx.Rollback()
		return model.Guru{}, "", err
	}

	if repository.CheckNIPExists(data.NIP, 0) {
		tx.Rollback()
		return model.Guru{}, "", errors.New("NIP sudah terdaftar")
	}

	// 3. Generate random temporary password
	plainPassword := GenerateRandomPassword()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return model.Guru{}, "", err
	}

	// 4. Create User record (Guru does NOT use force change password, so is_first_login = false)
	user := authModel.User{
		Nama:         data.Nama,
		Email:        &email,
		Password:     string(hashedPassword),
		Role:         "guru",
		IsFirstLogin: false,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return model.Guru{}, "", errors.New("gagal membuat akun user untuk guru: " + err.Error())
	}

	// 5. Link User ID to Guru record and save Guru
	data.UserID = user.ID
	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return model.Guru{}, "", errors.New("gagal menyimpan data guru: " + err.Error())
	}

	// Commit Transaction
	if err := tx.Commit().Error; err != nil {
		return model.Guru{}, "", err
	}

	config.DB.Preload("User").First(&data, data.ID)

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
		strings.TrimSpace(data.JenisKelamin) == "" {
		return errors.New("nama, gelar, dan jenis_kelamin wajib diisi")
	}

	return nil
}
