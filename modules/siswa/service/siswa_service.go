package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"backend/config"
	academicModel "backend/modules/academic/model"
	academicService "backend/modules/academic/service"
	authModel "backend/modules/auth/model"
	"backend/modules/siswa/model"
	"backend/modules/siswa/repository"

	"golang.org/x/crypto/bcrypt"
)

func GetAllSiswa(ctx academicModel.TeacherContext, search string, kelasID string) ([]model.Siswa, error) {
	if ctx.Role == "guru" {
		// If teacher, only return students taught by this teacher.
		// If kelasID is provided, filter by that as well.
		if kelasID != "" {
			parsedKelas, _ := strconv.Atoi(kelasID)
			return academicService.GetStudentsByTeacherAndClass(ctx, uint(parsedKelas))
		}
		return academicService.GetStudentsByTeacher(ctx)
	}

	// For admin/other roles
	return repository.GetAllSiswa(search, kelasID, "")
}

func GetSiswaByID(ctx academicModel.TeacherContext, id uint) (model.Siswa, error) {
	if ctx.Role == "guru" {
		if err := academicService.CanViewStudent(ctx, id); err != nil {
			return model.Siswa{}, err
		}
	}
	return repository.GetSiswaByID(id)
}

func CreateSiswa(data model.Siswa) (model.Siswa, string, error) {
	tx := repository.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	nis, err := repository.GenerateNISWithTx(tx)
	if err != nil {
		tx.Rollback()
		return model.Siswa{}, "", err
	}
	data.NIS = nis

	if err := validateSiswa(data); err != nil {
		tx.Rollback()
		return model.Siswa{}, "", err
	}

	var seq int
	if len(data.NIS) >= 4 {
		seqStr := data.NIS[len(data.NIS)-4:]
		seq, _ = strconv.Atoi(seqStr)
	}
	if seq == 0 {
		seq = 1
	}
	email := fmt.Sprintf("siswa%03d@sekolah.com", seq)

	exists, err := repository.CountSiswaByNISWithTx(tx, data.NIS)
	if err != nil {
		tx.Rollback()
		return model.Siswa{}, "", fmt.Errorf("ERR_GENERATE_NIS: Gagal memeriksa keunikan NIS: %v", err)
	}
	if exists > 0 {
		tx.Rollback()
		return model.Siswa{}, "", errors.New("ERR_IDENTIFIER_DUPLICATE: NIS sudah terdaftar")
	}

	emailCount, err := repository.CountUserByEmailWithTx(tx, email)
	if err != nil {
		tx.Rollback()
		return model.Siswa{}, "", fmt.Errorf("ERR_GENERATE_NIS: Gagal memeriksa keunikan email: %v", err)
	}
	if emailCount > 0 {
		tx.Rollback()
		return model.Siswa{}, "", errors.New("ERR_EMAIL_DUPLICATE: Email sudah terdaftar")
	}

	plainPassword := "Siswa123!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return model.Siswa{}, "", err
	}

	user := authModel.User{
		Nama:         data.Nama,
		Email:        &email,
		Password:     string(hashedPassword),
		Role:         "siswa",
		IsFirstLogin: true,
	}

	if err := repository.CreateUserWithTx(tx, &user); err != nil {
		tx.Rollback()
		return model.Siswa{}, "", errors.New("ERR_EMAIL_DUPLICATE: gagal membuat akun user untuk siswa: " + err.Error())
	}

	data.UserID = user.ID
	if err := repository.CreateSiswaWithTx(tx, data); err != nil {
		tx.Rollback()
		return model.Siswa{}, "", errors.New("ERR_GENERATE_NIS: gagal menyimpan data siswa: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return model.Siswa{}, "", err
	}

	if config.Debug {
		log.Printf("[INFO] Action: Student Created | Timestamp: %s | User: %s | Identifier: %s | Email: %s", time.Now().Format(time.RFC3339), data.Nama, data.NIS, email)
	}

	repository.LoadSiswaRelations(&data)
	return data, plainPassword, nil
}

func UpdateSiswa(id uint, data model.Siswa) (model.Siswa, error) {
	if err := validateSiswa(data); err != nil {
		return model.Siswa{}, err
	}
	return repository.UpdateSiswa(id, data)
}

func DeleteSiswa(id uint) error {
	return repository.DeleteSiswa(id)
}

func validateSiswa(data model.Siswa) error {
	if data.Nama == "" || data.JenisKelamin == "" || data.TanggalLahir == "" ||
		data.Provinsi == "" || data.Kabupaten == "" || data.Kecamatan == "" ||
		data.Desa == "" || data.AlamatDetail == "" || data.KelasID == 0 {
		return errors.New("nama, jenis_kelamin, tanggal_lahir, provinsi, kabupaten, kecamatan, desa, alamat_detail, dan kelas_id wajib diisi")
	}

	kelas, err := repository.GetKelasByID(data.KelasID)
	if err != nil {
		return errors.New("kelas tidak ditemukan")
	}

	count, err := repository.CountSiswaByKelasExcludeID(data.KelasID, data.ID)
	if err != nil {
		return err
	}

	if int(count) >= kelas.Kapasitas {
		return errors.New("Kelas sudah mencapai kapasitas maksimum.")
	}

	return nil
}
