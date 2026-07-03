package model

import (
	"time"

	guruModel "backend/modules/guru/model"
	siswaModel "backend/modules/siswa/model"

	"gorm.io/gorm"
)

type Perizinan struct {
	ID             uint             `gorm:"primaryKey" json:"id"`
	SiswaID        uint             `gorm:"not null" json:"siswa_id"`
	Siswa          siswaModel.Siswa `gorm:"foreignKey:SiswaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"siswa,omitempty"`
	TanggalMulai   time.Time        `gorm:"type:date;not null" json:"tanggal_mulai"`
	TanggalSelesai time.Time        `gorm:"type:date;not null" json:"tanggal_selesai"`
	Tipe           string           `gorm:"type:varchar(20);not null" json:"tipe"` // "Izin" atau "Sakit"
	Alasan         string           `gorm:"type:text;not null" json:"alasan"`
	AttachmentURL  *string          `gorm:"type:varchar(255)" json:"attachment_url"`
	Status         string           `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"` // "Pending", "Disetujui", "Ditolak"
	KeteranganGuru *string          `gorm:"type:text" json:"keterangan_guru"`
	DisetujuiOleh  *uint            `json:"disetujui_oleh"`
	Guru           *guruModel.Guru  `gorm:"foreignKey:DisetujuiOleh;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"guru,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	DeletedAt      gorm.DeletedAt   `gorm:"index" json:"-"`
}

func (Perizinan) TableName() string {
	return "perizinan"
}
