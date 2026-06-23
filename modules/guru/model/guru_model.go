package model

import (
	"time"

	"gorm.io/gorm"
)

type Guru struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	NIP          string         `gorm:"type:varchar(30);uniqueIndex;not null" json:"nip"`
	Nama         string         `gorm:"type:varchar(100);not null" json:"nama"`
	JenisKelamin string         `gorm:"type:varchar(20);not null" json:"jenis_kelamin"`
	Email        string         `gorm:"type:varchar(120)" json:"email"`
	NoHP         string         `gorm:"type:varchar(30)" json:"no_hp"`
	Alamat       string         `gorm:"type:text" json:"alamat"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
