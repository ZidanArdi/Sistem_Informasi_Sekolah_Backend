package model

import (
	"time"

	guruModel "backend/modules/guru/model"
	siswaModel "backend/modules/siswa/model"

	"gorm.io/gorm"
)

type Perizinan struct {
	ID             uint             `gorm:"primaryKey" json:"id" example:"1"`
	SiswaID        uint             `gorm:"not null" json:"siswa_id" example:"1"`
	Siswa          siswaModel.Siswa `gorm:"foreignKey:SiswaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"siswa,omitempty"`
	TanggalMulai   time.Time        `gorm:"type:date;not null" json:"tanggal_mulai" example:"2026-07-15T00:00:00Z"`
	TanggalSelesai time.Time        `gorm:"type:date;not null" json:"tanggal_selesai" example:"2026-07-16T00:00:00Z"`
	Tipe           string           `gorm:"type:varchar(20);not null" json:"tipe" example:"Izin"` // "Izin" atau "Sakit"
	Alasan         string           `gorm:"type:text;not null" json:"alasan" example:"Ada keperluan keluarga"`
	AttachmentURL  *string          `gorm:"type:varchar(255)" json:"attachment_url" example:"/uploads/perizinan/sample.pdf"`
	Status         string           `gorm:"type:varchar(20);not null;default:'Pending'" json:"status" example:"Pending"` // "Pending", "Disetujui", "Ditolak"
	KeteranganGuru *string          `gorm:"type:text" json:"keterangan_guru" example:"Semoga lekas sembuh/urusan lancar"`
	DisetujuiOleh  *uint            `json:"disetujui_oleh" example:"1"`
	Guru           *guruModel.Guru  `gorm:"foreignKey:DisetujuiOleh;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"guru,omitempty"`
	CreatedAt      time.Time        `json:"created_at" example:"2026-07-15T19:00:00Z"`
	UpdatedAt      time.Time        `json:"updated_at" example:"2026-07-15T19:00:00Z"`
	DeletedAt      gorm.DeletedAt   `gorm:"index" json:"-"`
}

func (Perizinan) TableName() string {
	return "perizinan"
}

// PerizinanResponseSingle represents a successful single Perizinan response
type PerizinanResponseSingle struct {
	Success bool      `json:"success" example:"true"`
	Message string    `json:"message" example:"Berhasil mengambil detail perizinan"`
	Data    Perizinan `json:"data"`
}

// PerizinanResponseList represents a successful list of Perizinan response
type PerizinanResponseList struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Berhasil mengambil data perizinan"`
	Data    []Perizinan `json:"data"`
}
