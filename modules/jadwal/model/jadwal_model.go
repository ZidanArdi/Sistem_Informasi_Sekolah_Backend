package model

import (
	"time"

	guruModel "backend/modules/guru/model"
	kelasModel "backend/modules/kelas/model"
	mapelModel "backend/modules/mapel/model"

	"gorm.io/gorm"
)

type Jadwal struct {
	ID          uint             `gorm:"primaryKey" json:"id" example:"1"`
	KelasID     uint             `gorm:"not null" json:"kelas_id" example:"1"`
	MapelID     uint             `gorm:"not null" json:"mapel_id" example:"1"`
	GuruID      uint             `gorm:"not null" json:"guru_id" example:"1"`
	Hari        string           `gorm:"type:varchar(20);not null" json:"hari" example:"Senin"`
	JamMulai    string           `gorm:"type:time;not null" json:"jam_mulai" example:"07:00"`
	JamSelesai  string           `gorm:"type:time;not null" json:"jam_selesai" example:"09:00"`
	TahunAjaran string           `gorm:"type:varchar(10);not null" json:"tahun_ajaran" example:"2026/2027"`
	Semester    string           `gorm:"type:varchar(10);not null" json:"semester" example:"Ganjil"`
	Kelas       kelasModel.Kelas `gorm:"foreignKey:KelasID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"kelas,omitempty"`
	Mapel       mapelModel.Mapel `gorm:"foreignKey:MapelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"mapel,omitempty"`
	Guru        guruModel.Guru   `gorm:"foreignKey:GuruID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"guru,omitempty"`
	CreatedAt   time.Time        `json:"created_at" example:"2026-07-15T19:00:00Z"`
	UpdatedAt   time.Time        `json:"updated_at" example:"2026-07-15T19:00:00Z"`
	DeletedAt   gorm.DeletedAt   `gorm:"index" json:"-"`
}

// JadwalResponseSingle represents a successful single Jadwal response
type JadwalResponseSingle struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Berhasil mengambil detail jadwal pelajaran"`
	Data    Jadwal `json:"data"`
}

// JadwalResponseList represents a successful list of Jadwal response
type JadwalResponseList struct {
	Success bool     `json:"success" example:"true"`
	Message string   `json:"message" example:"Berhasil mengambil data jadwal pelajaran"`
	Data    []Jadwal `json:"data"`
}
