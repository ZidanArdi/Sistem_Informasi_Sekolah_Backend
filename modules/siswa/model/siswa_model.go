package model

import (
	"time"

	kelasModel "backend/modules/kelas/model"

	"gorm.io/gorm"
)

type Siswa struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	UserID       uint             `gorm:"index" json:"user_id"`
	NIS          string           `gorm:"type:varchar(30);uniqueIndex;not null" json:"nis"`
	Nama         string           `gorm:"type:varchar(100);not null" json:"nama"`
	JenisKelamin string           `gorm:"type:varchar(20)" json:"jenis_kelamin"`
	TanggalLahir string           `gorm:"type:date" json:"tanggal_lahir"`
	NoHP         string           `gorm:"type:varchar(20)" json:"no_hp"`
	PhotoURL     string           `gorm:"type:varchar(255)" json:"photo_url"`
	KelasID      uint             `json:"kelas_id"`
	Kelas        kelasModel.Kelas `gorm:"foreignKey:KelasID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"kelas,omitempty"`
	
	// Address Normalization Fields
	Provinsi     string `gorm:"type:varchar(100)" json:"provinsi"`
	Kabupaten    string `gorm:"type:varchar(100)" json:"kabupaten"`
	Kecamatan    string `gorm:"type:varchar(100)" json:"kecamatan"`
	Desa         string `gorm:"type:varchar(100)" json:"desa"`
	AlamatDetail string `gorm:"type:text" json:"alamat_detail"`

	// Deprecated Database Columns (Retained for Non-Destructive Migrations)
	TempatLahir string  `gorm:"type:varchar(80)" json:"tempat_lahir,omitempty"` // Deprecated
	Alamat      string  `gorm:"type:text" json:"alamat,omitempty"`             // Deprecated
	Email       *string `gorm:"type:varchar(120);uniqueIndex" json:"email,omitempty"` // Deprecated

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
