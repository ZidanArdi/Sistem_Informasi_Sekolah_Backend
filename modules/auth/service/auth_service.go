package service

import (
	"errors"
	"strings"

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
		return model.User{}, errors.New("email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	emailVal := strings.TrimSpace(strings.ToLower(input.Email))
	user := model.User{
		Nama:     strings.TrimSpace(input.Nama),
		Email:    &emailVal,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	return repository.CreateUser(user)
}

func Login(input LoginInput) (model.User, string, error) {
	loginID := input.Identifier
	if loginID == "" {
		loginID = input.Email
	}

	if strings.TrimSpace(loginID) == "" || input.Password == "" {
		return model.User{}, "", errors.New("email/NIS dan password wajib diisi")
	}

	loginID = strings.TrimSpace(loginID)
	var user model.User
	var err error

	if strings.Contains(loginID, "@") {
		// Email Login (Admin/Guru)
		user, err = repository.GetUserByEmail(strings.ToLower(loginID))
		if err != nil {
			return model.User{}, "", errors.New("email atau password salah")
		}
		if user.Role == "siswa" {
			return model.User{}, "", errors.New("siswa wajib login menggunakan NIS")
		}
	} else {
		// NIS Login (Siswa)
		userID, _, errNIS := repository.GetUserIDByNIS(loginID)
		if errNIS != nil {
			return model.User{}, "", errors.New("NIS tidak terdaftar")
		}
		if userID == 0 {
			return model.User{}, "", errors.New("akun siswa belum terbuat, silakan hubungi admin")
		}

		user, err = repository.GetUserByID(userID)
		if err != nil {
			return model.User{}, "", errors.New("akun siswa tidak ditemukan")
		}
		if user.Role != "siswa" {
			return model.User{}, "", errors.New("NIS hanya untuk login siswa")
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return model.User{}, "", errors.New("email/NIS atau password salah")
	}

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
			Alamat       string
			JenisKelamin string
		}
		config.DB.Table("siswas").
			Select("nis, no_hp, alamat, jenis_kelamin").
			Where("user_id = ?", user.ID).
			Scan(&siswa)
		user.NIS = siswa.NIS
		user.NoHP = siswa.NoHP
		user.Alamat = siswa.Alamat
		user.JenisKelamin = siswa.JenisKelamin
	} else if user.Role == "guru" {
		var guru struct {
			NIP          string
			Gelar        string
			NoHP         string
			Alamat       string
			JenisKelamin string
		}
		config.DB.Table("gurus").
			Select("nip, gelar, no_hp, alamat, jenis_kelamin").
			Where("user_id = ?", user.ID).
			Scan(&guru)
		user.NIP = guru.NIP
		user.Gelar = guru.Gelar
		user.NoHP = guru.NoHP
		user.Alamat = guru.Alamat
		user.JenisKelamin = guru.JenisKelamin
	}

	return user, token, nil
}

func ChangePassword(userID uint, input ChangePasswordInput) error {
	if input.OldPassword == "" || input.NewPassword == "" {
		return errors.New("password lama dan password baru wajib diisi")
	}

	if len(input.NewPassword) < 6 {
		return errors.New("password baru minimal 6 karakter")
	}

	user, err := repository.GetUserByID(userID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		return errors.New("password lama tidak sesuai")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return repository.UpdatePassword(userID, string(hashedPassword))
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

	if input.Role != "admin" && input.Role != "guru" {
		return errors.New("pendaftaran mandiri hanya diperbolehkan untuk role admin atau guru")
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
