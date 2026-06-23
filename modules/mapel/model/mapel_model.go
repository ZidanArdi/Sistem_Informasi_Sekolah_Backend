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
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
