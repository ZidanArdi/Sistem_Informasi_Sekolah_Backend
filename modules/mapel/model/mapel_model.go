package model

import (
	"time"

	"gorm.io/gorm"
)

type Mapel struct {
	ID        uint           `gorm:"primaryKey" json:"id" example:"1"`
	KodeMapel string         `gorm:"type:varchar(30);uniqueIndex;not null" json:"kode_mapel" example:"MP001"`
	NamaMapel string         `gorm:"type:varchar(100);not null" json:"nama_mapel" example:"Matematika"`
	Jam       int            `gorm:"not null" json:"jam" example:"4"`
	Jurusan   string         `gorm:"type:varchar(50)" json:"jurusan" example:"Umum"`
	IsUmum    bool           `gorm:"not null;default:true" json:"is_umum" example:"true"`
	CreatedAt time.Time      `json:"created_at" example:"2026-07-15T19:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2026-07-15T19:00:00Z"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// MapelResponseSingle represents a successful single Mapel response
type MapelResponseSingle struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Berhasil mengambil detail mata pelajaran"`
	Data    Mapel  `json:"data"`
}

// MapelResponseList represents a successful list of Mapel response
type MapelResponseList struct {
	Success bool    `json:"success" example:"true"`
	Message string  `json:"message" example:"Berhasil mengambil data mata pelajaran"`
	Data    []Mapel `json:"data"`
}
