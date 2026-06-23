package service

import (
	"errors"
	"strings"

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
	Email    string `json:"email"`
	Password string `json:"password"`
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

	user := model.User{
		Nama:     strings.TrimSpace(input.Nama),
		Email:    strings.TrimSpace(strings.ToLower(input.Email)),
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	return repository.CreateUser(user)
}

func Login(input LoginInput) (model.User, string, error) {
	if strings.TrimSpace(input.Email) == "" || input.Password == "" {
		return model.User{}, "", errors.New("email dan password wajib diisi")
	}

	user, err := repository.GetUserByEmail(strings.TrimSpace(strings.ToLower(input.Email)))
	if err != nil {
		return model.User{}, "", errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return model.User{}, "", errors.New("email atau password salah")
	}

	token, err := helpers.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return model.User{}, "", err
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

	if input.Role != "admin" && input.Role != "guru" && input.Role != "staff" && input.Role != "siswa" {
		return errors.New("role harus admin, guru, staff, atau siswa")
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
