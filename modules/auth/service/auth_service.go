package service

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"backend/helpers"
	"backend/modules/auth/model"
	"backend/modules/auth/repository"

	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Nama     string `json:"nama" example:"Ahmad Guru"`
	Email    string `json:"email" example:"ahmad.guru@sekolah.com"`
	Password string `json:"password" example:"Guru123!"`
	Role     string `json:"role" example:"guru"`
}

type LoginInput struct {
	Email      string `json:"email" example:"admin@sekolah.com"`
	Identifier string `json:"identifier" example:"ADM001"`
	Password   string `json:"password" example:"Admin123!"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" example:"Guru123!"`
	NewPassword string `json:"new_password" example:"GuruNew123!"`
}

func Register(input RegisterInput) (model.User, error) {
	input.Role = normalizeRole(input.Role)
	if err := validateRegister(input); err != nil {
		return model.User{}, err
	}

	if repository.EmailExists(input.Email) {
		return model.User{}, errors.New("ERR_EMAIL_DUPLICATE: email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	emailVal := strings.TrimSpace(strings.ToLower(input.Email))
	user := model.User{
		Nama:         strings.TrimSpace(input.Nama),
		Email:        &emailVal,
		Password:     string(hashedPassword),
		Role:         input.Role,
		IsFirstLogin: true, // Every newly created account IsFirstLogin = true
	}

	return repository.CreateUser(user)
}

func Login(input LoginInput) (model.User, string, error) {
	loginID := input.Identifier
	if loginID == "" {
		loginID = input.Email
	}

	if strings.TrimSpace(loginID) == "" || input.Password == "" {
		return model.User{}, "", errors.New("ERR_LOGIN_FAILED: email/ID Pengguna dan password wajib diisi")
	}

	loginID = strings.TrimSpace(loginID)
	user, err := repository.FindLoginIdentity(loginID)
	if err != nil {
		return model.User{}, "", err
	}

	if !user.IsActive {
		return model.User{}, "", errors.New("ERR_PERMISSION_DENIED: akun Anda telah dinonaktifkan, silakan hubungi admin")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return model.User{}, "", errors.New("ERR_LOGIN_FAILED: email/ID Pengguna atau password salah")
	}

	// First login flag is fetched directly from database User.IsFirstLogin.
	// It will be updated to false when the user changes their password via ChangePassword.


	// Update last login timestamp
	now := time.Now()
	user.LastLoginAt = &now
	repository.UpdateLastLogin(user.ID, &now)

	var emailStr string
	if user.Email != nil {
		emailStr = *user.Email
	} else {
		emailStr = loginID // Use NIS as placeholder
	}

	token, err := helpers.GenerateToken(user.ID, emailStr, user.Role)
	if err != nil {
		return model.User{}, "", err
	}

	// Populate profile details based on role
	if user.Role == "siswa" {
		siswa, _ := repository.GetSiswaProfileByUserID(user.ID)
		user.NIS = siswa.NIS
		user.NoHP = siswa.NoHP
		user.JenisKelamin = siswa.JenisKelamin
		if siswa.AlamatDetail != "" {
			user.Alamat = fmt.Sprintf("%s, %s, %s, %s, %s", siswa.AlamatDetail, siswa.Desa, siswa.Kecamatan, siswa.Kabupaten, siswa.Provinsi)
		}
	} else if user.Role == "guru" {
		guru, _ := repository.GetGuruProfileByUserID(user.ID)
		user.NIP = guru.NIP
		user.Gelar = guru.Gelar
		user.NoHP = guru.NoHP
		user.JenisKelamin = guru.JenisKelamin
		if guru.AlamatDetail != "" {
			user.Alamat = fmt.Sprintf("%s, %s, %s, %s, %s", guru.AlamatDetail, guru.Desa, guru.Kecamatan, guru.Kabupaten, guru.Provinsi)
		}
	}

	// LOG Login Info
	var resolvedIdentifier string
	if user.Role == "admin" {
		resolvedIdentifier = user.Nama
	} else if user.Role == "guru" {
		nip, _ := repository.GetGuruNIPByUserID(user.ID)
		resolvedIdentifier = nip
	} else {
		nis, _ := repository.GetSiswaNISByUserID(user.ID)
		resolvedIdentifier = nis
	}
	log.Printf("[INFO] Action: Login | Timestamp: %s | User: %s | Identifier: %s | Email: %s", time.Now().Format(time.RFC3339), user.Nama, resolvedIdentifier, emailStr)

	return user, token, nil
}

func validateNewPasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password baru minimal 8 karakter")
	}
	var hasUpper, hasLower, hasNumber bool
	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		} else if char >= 'a' && char <= 'z' {
			hasLower = true
		} else if char >= '0' && char <= '9' {
			hasNumber = true
		}
	}
	if !hasUpper {
		return errors.New("password baru harus mengandung minimal 1 huruf besar")
	}
	if !hasLower {
		return errors.New("password baru harus mengandung minimal 1 huruf kecil")
	}
	if !hasNumber {
		return errors.New("password baru harus mengandung minimal 1 angka")
	}
	return nil
}

func ChangePassword(userID uint, input ChangePasswordInput) error {
	if input.OldPassword == "" || input.NewPassword == "" {
		return errors.New("password lama dan password baru wajib diisi")
	}

	if err := validateNewPasswordStrength(input.NewPassword); err != nil {
		return err
	}

	user, err := repository.GetUserByID(userID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		return errors.New("ERR_LOGIN_FAILED: password lama tidak sesuai")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = repository.UpdatePassword(userID, string(hashedPassword))
	if err != nil {
		return err
	}

	// LOG Password Changed Info
	var resolvedIdentifier string
	var emailStr string
	if user.Email != nil {
		emailStr = *user.Email
	}

	if user.Role == "admin" {
		resolvedIdentifier = user.Nama
	} else if user.Role == "guru" {
		nip, _ := repository.GetGuruNIPByUserID(user.ID)
		resolvedIdentifier = nip
	} else {
		nis, _ := repository.GetSiswaNISByUserID(user.ID)
		resolvedIdentifier = nis
	}
	log.Printf("[INFO] Action: Password Changed | Timestamp: %s | User: %s | Identifier: %s | Email: %s", time.Now().Format(time.RFC3339), user.Nama, resolvedIdentifier, emailStr)

	return nil
}

func validateRegister(input RegisterInput) error {
	if strings.TrimSpace(input.Nama) == "" ||
		strings.TrimSpace(input.Email) == "" ||
		input.Password == "" {
		return errors.New("nama, email, dan password wajib diisi")
	}

	if !strings.Contains(input.Email, "@") {
		return errors.New("format email tidak valid")
	}

	if len(input.Password) < 6 {
		return errors.New("password minimal 6 karakter")
	}

	if input.Role != "guru" {
		return errors.New("ERR_PERMISSION_DENIED: pendaftaran mandiri hanya diperbolehkan untuk role guru")
	}

	return nil
}

func normalizeRole(role string) string {
	role = strings.TrimSpace(strings.ToLower(role))
	if role == "" {
		return "staff"
	}
	return role
}
