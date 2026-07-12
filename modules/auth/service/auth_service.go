package service

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"backend/config"
	"backend/helpers"
	"backend/modules/auth/model"
	"backend/modules/auth/repository"

	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginInput struct {
	Email      string `json:"email"`
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
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
	var user model.User
	var err error

	if strings.Contains(loginID, "@") {
		// Email Login
		user, err = repository.GetUserByEmail(strings.ToLower(loginID))
		if err != nil {
			return model.User{}, "", errors.New("ERR_INVALID_IDENTIFIER: email/ID Pengguna atau password salah")
		}
	} else if strings.HasPrefix(strings.ToUpper(loginID), "GR") {
		// NIG Login (Guru)
		var guru struct {
			UserID uint `gorm:"column:user_id"`
		}
		errNIG := config.DB.Table("gurus").Select("user_id").Where("nip = ? AND deleted_at IS NULL", strings.ToUpper(loginID)).First(&guru).Error
		if errNIG != nil {
			return model.User{}, "", errors.New("ERR_INVALID_IDENTIFIER: NIG tidak terdaftar")
		}
		user, err = repository.GetUserByID(guru.UserID)
		if err != nil {
			return model.User{}, "", errors.New("ERR_INVALID_IDENTIFIER: akun guru tidak ditemukan")
		}
	} else if strings.ToUpper(loginID) == "ADM001" {
		// Admin Code Login
		user, err = repository.GetUserByEmail("admin@sekolah.com")
		if err != nil {
			return model.User{}, "", errors.New("ERR_INVALID_IDENTIFIER: akun admin tidak ditemukan")
		}
	} else {
		// NIS Login (Siswa)
		userID, _, errNIS := repository.GetUserIDByNIS(loginID)
		if errNIS != nil {
			return model.User{}, "", errors.New("ERR_INVALID_IDENTIFIER: NIS tidak terdaftar")
		}
		user, err = repository.GetUserByID(userID)
		if err != nil {
			return model.User{}, "", errors.New("ERR_INVALID_IDENTIFIER: akun siswa tidak ditemukan")
		}
	}

	if !user.IsActive {
		return model.User{}, "", errors.New("ERR_PERMISSION_DENIED: akun Anda telah dinonaktifkan, silakan hubungi admin")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return model.User{}, "", errors.New("ERR_LOGIN_FAILED: email/ID Pengguna atau password salah")
	}

	// Dynamic first login / default password check
	isDefaultPassword := false
	if user.Role == "admin" && input.Password == "Admin123!" {
		isDefaultPassword = true
	} else if user.Role == "guru" && input.Password == "Guru123!" {
		isDefaultPassword = true
	} else if user.Role == "siswa" && input.Password == "Siswa123!" {
		isDefaultPassword = true
	}

	if isDefaultPassword {
		user.IsFirstLogin = true
		config.DB.Model(&model.User{}).Where("id = ?", user.ID).Update("is_first_login", true)
	}

	// Update last login timestamp
	now := time.Now()
	user.LastLoginAt = &now
	config.DB.Model(&model.User{}).Where("id = ?", user.ID).Update("last_login_at", &now)

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
		var siswa struct {
			NIS          string
			NoHP         string
			JenisKelamin string
			Provinsi     string
			Kabupaten    string
			Kecamatan    string
			Desa         string
			AlamatDetail string
		}
		config.DB.Table("siswas").
			Select("nis, no_hp, jenis_kelamin, provinsi, kabupaten, kecamatan, desa, alamat_detail").
			Where("user_id = ?", user.ID).
			Scan(&siswa)
		user.NIS = siswa.NIS
		user.NoHP = siswa.NoHP
		user.JenisKelamin = siswa.JenisKelamin
		if siswa.AlamatDetail != "" {
			user.Alamat = fmt.Sprintf("%s, %s, %s, %s, %s", siswa.AlamatDetail, siswa.Desa, siswa.Kecamatan, siswa.Kabupaten, siswa.Provinsi)
		}
	} else if user.Role == "guru" {
		var guru struct {
			NIP          string
			Gelar        string
			NoHP         string
			JenisKelamin string
			Provinsi     string
			Kabupaten    string
			Kecamatan    string
			Desa         string
			AlamatDetail string
		}
		config.DB.Table("gurus").
			Select("nip, gelar, no_hp, jenis_kelamin, provinsi, kabupaten, kecamatan, desa, alamat_detail").
			Where("user_id = ?", user.ID).
			Scan(&guru)
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
		resolvedIdentifier = "ADM001"
	} else if user.Role == "guru" {
		resolvedIdentifier = user.NIP
	} else {
		resolvedIdentifier = user.NIS
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
		resolvedIdentifier = "ADM001"
	} else if user.Role == "guru" {
		var nip string
		config.DB.Table("gurus").Select("nip").Where("user_id = ?", user.ID).Row().Scan(&nip)
		resolvedIdentifier = nip
	} else {
		var nis string
		config.DB.Table("siswas").Select("nis").Where("user_id = ?", user.ID).Row().Scan(&nis)
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
