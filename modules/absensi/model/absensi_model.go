package model

import (
	"time"

	siswaModel "backend/modules/siswa/model"

	"gorm.io/gorm"
)

type Absensi struct {
	ID                uint             `gorm:"primaryKey" json:"id"`
	SiswaID           uint             `gorm:"not null" json:"siswa_id"`
	Siswa             siswaModel.Siswa `gorm:"foreignKey:SiswaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"siswa,omitempty"`
	Tanggal           string           `gorm:"type:date;not null" json:"tanggal"` // format: YYYY-MM-DD
	Status            string           `gorm:"type:varchar(20);not null" json:"status"` // "Hadir", "Sakit", "Izin", "Alpa"
	Keterangan        string           `gorm:"type:text" json:"keterangan"` // detail alasan jika Sakit/Izin
	StatusPersetujuan string           `gorm:"type:varchar(50);not null;default:'Disetujui'" json:"status_persetujuan"` // "Pending", "Disetujui", "Ditolak"
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         gorm.DeletedAt   `gorm:"index" json:"-"`
}

func (Absensi) TableName() string {
	return "absensi"
}
