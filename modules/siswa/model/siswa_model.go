package model

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Siswa struct {
	NIS       string         `gorm:"primaryKey;type:varchar(20)" json:"nis"`
	Nama      string         `json:"nama"`
	Kelas     string         `json:"kelas"`
	Alamat    string         `json:"alamat"`
	Email     string         `gorm:"unique" json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Siswa) BeforeCreate(tx *gorm.DB) (err error) {
	var lastSiswa Siswa
	// Cari siswa terakhir berdasarkan NIS (order by NIS descending)
	result := tx.Order("nis desc").First(&lastSiswa)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Jika belum ada data sama sekali, berikan NIS awal
			s.NIS = "SIS001"
			return nil
		}
		return result.Error
	}

	// Ambil angka dari NIS terakhir, misalnya "SIS001" -> "001" -> 1
	if len(lastSiswa.NIS) > 3 {
		lastNumberStr := lastSiswa.NIS[3:]
		lastNumber, err := strconv.Atoi(lastNumberStr)
		if err != nil {
			return err
		}
		// Tambahkan 1 dan format kembali menjadi 3 digit (SIS002, SIS003, dst)
		s.NIS = fmt.Sprintf("SIS%03d", lastNumber+1)
	} else {
		s.NIS = "SIS001"
	}

	return nil
}