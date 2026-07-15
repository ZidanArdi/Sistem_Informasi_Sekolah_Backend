package model

import (
	"time"

	guruModel "backend/modules/guru/model"

	"gorm.io/gorm"
)

type Kelas struct {
	ID          uint            `gorm:"primaryKey" json:"id" example:"1"`
	NamaKelas   string          `gorm:"type:varchar(80);not null" json:"nama_kelas" example:"X RPL 1"`
	Tingkat     string          `gorm:"type:varchar(20);not null" json:"tingkat" example:"X"`
	Jurusan     string          `gorm:"type:varchar(50)" json:"jurusan" example:"Rekayasa Perangkat Lunak"`
	Kapasitas   int             `gorm:"not null;default:0" json:"kapasitas" example:"36"`
	WaliKelasID *uint           `gorm:"uniqueIndex" json:"wali_kelas_id" example:"1"`
	WaliKelas   *guruModel.Guru `gorm:"foreignKey:WaliKelasID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"wali_kelas,omitempty"`
	TotalSiswa  int             `gorm:"-" json:"total_siswa" example:"32"`
	CreatedAt   time.Time       `json:"created_at" example:"2026-07-15T19:00:00Z"`
	UpdatedAt   time.Time       `json:"updated_at" example:"2026-07-15T19:00:00Z"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

// KelasResponseSingle represents a successful single Kelas response
type KelasResponseSingle struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Berhasil mengambil detail kelas"`
	Data    Kelas  `json:"data"`
}

// KelasResponseList represents a successful list of Kelas response
type KelasResponseList struct {
	Success bool    `json:"success" example:"true"`
	Message string  `json:"message" example:"Berhasil mengambil data kelas"`
	Data    []Kelas `json:"data"`
}
