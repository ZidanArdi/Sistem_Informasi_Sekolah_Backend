package model

import (
	"time"

	authModel "backend/modules/auth/model"

	"gorm.io/gorm"
)

type Guru struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	UserID       uint            `gorm:"index" json:"user_id"`
	User         *authModel.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	NIP          string          `gorm:"column:nip;type:varchar(30);uniqueIndex;not null" json:"nip"`
	Nama         string          `gorm:"type:varchar(100);not null" json:"nama"`
	Gelar        string          `gorm:"type:varchar(50);not null" json:"gelar"`
	JenisKelamin string          `gorm:"type:varchar(20);not null" json:"jenis_kelamin"`
	NoHP         string          `gorm:"type:varchar(30)" json:"no_hp"`
	PhotoURL     string          `gorm:"type:varchar(255)" json:"photo_url"`
	
	// Address Normalization Fields
	Provinsi     string `gorm:"type:varchar(100)" json:"provinsi"`
	Kabupaten    string `gorm:"type:varchar(100)" json:"kabupaten"`
	Kecamatan    string `gorm:"type:varchar(100)" json:"kecamatan"`
	Desa         string `gorm:"type:varchar(100)" json:"desa"`
	AlamatDetail string `gorm:"type:text" json:"alamat_detail"`

	// Deprecated Database Columns (Retained for Non-Destructive Migrations)
	Alamat string `gorm:"type:text" json:"alamat,omitempty"` // Deprecated

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Virtual field for frontend multiselect checkboxes
	MapelIDs  []uint         `gorm:"-" json:"mapel_ids,omitempty"`
}

type GuruMapel struct {
	ID      uint `gorm:"primaryKey" json:"id"`
	GuruID  uint `gorm:"uniqueIndex:idx_guru_mapel;not null" json:"guru_id"`
	MapelID uint `gorm:"uniqueIndex:idx_guru_mapel;not null" json:"mapel_id"`
}

func (GuruMapel) TableName() string {
	return "guru_mapel"
}
