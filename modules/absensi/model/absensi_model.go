package model

import (
	"time"

	siswaModel "backend/modules/siswa/model"

	"gorm.io/gorm"
)

type Absensi struct {
	ID                uint             `gorm:"primaryKey" json:"id" example:"1"`
	SiswaID           uint             `gorm:"not null" json:"siswa_id" example:"1"`
	Siswa             siswaModel.Siswa `gorm:"foreignKey:SiswaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"siswa,omitempty"`
	Tanggal           string           `gorm:"type:date;not null" json:"tanggal" example:"2026-07-15"` // format: YYYY-MM-DD
	Status            string           `gorm:"type:varchar(20);not null" json:"status" example:"Hadir"` // "Hadir", "Sakit", "Izin", "Alpa"
	Keterangan        string           `gorm:"type:text" json:"keterangan" example:"Hadir tepat waktu"` // detail alasan jika Sakit/Izin
	StatusPersetujuan string           `gorm:"type:varchar(50);not null;default:'Disetujui'" json:"status_persetujuan" example:"Disetujui"` // "Pending", "Disetujui", "Ditolak"
	CreatedAt         time.Time        `json:"created_at" example:"2026-07-15T19:00:00Z"`
	UpdatedAt         time.Time        `json:"updated_at" example:"2026-07-15T19:00:00Z"`
	DeletedAt         gorm.DeletedAt   `gorm:"index" json:"-"`
}

func (Absensi) TableName() string {
	return "absensi"
}

// AbsensiResponseSingle represents a successful single Absensi response
type AbsensiResponseSingle struct {
	Success bool    `json:"success" example:"true"`
	Message string  `json:"message" example:"Berhasil mengambil detail absensi"`
	Data    Absensi `json:"data"`
}

// AbsensiResponseList represents a successful list of Absensi response
type AbsensiResponseList struct {
	Success bool      `json:"success" example:"true"`
	Message string    `json:"message" example:"Berhasil mengambil data absensi"`
	Data    []Absensi `json:"data"`
}
