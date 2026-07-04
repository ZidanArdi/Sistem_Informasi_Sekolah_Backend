package model

import (
	"time"

	"gorm.io/gorm"
)

type Mapel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	KodeMapel string         `gorm:"type:varchar(30);uniqueIndex;not null" json:"kode_mapel"`
	NamaMapel string         `gorm:"type:varchar(100);not null" json:"nama_mapel"`
	Jam       int            `gorm:"not null" json:"jam"`
	Jurusan   string         `gorm:"type:varchar(50)" json:"jurusan"`
	IsUmum    bool           `gorm:"not null;default:true" json:"is_umum"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
