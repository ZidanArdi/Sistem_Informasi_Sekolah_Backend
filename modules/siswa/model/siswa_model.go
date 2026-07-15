package model

import (
	"time"

	kelasModel "backend/modules/kelas/model"

	"gorm.io/gorm"
)

type Siswa struct {
	ID           uint             `gorm:"primaryKey" json:"id" example:"1"`
	UserID       uint             `gorm:"index" json:"user_id" example:"2"`
	NIS          string           `gorm:"type:varchar(30);uniqueIndex;not null" json:"nis" example:"12345678"`
	Nama         string           `gorm:"type:varchar(100);not null" json:"nama" example:"Budi Santoso"`
	JenisKelamin string           `gorm:"type:varchar(20)" json:"jenis_kelamin" example:"Laki-laki"`
	TanggalLahir string           `gorm:"type:date" json:"tanggal_lahir" example:"2010-05-15"`
	NoHP         string           `gorm:"type:varchar(20)" json:"no_hp" example:"081234567890"`
	PhotoURL     string           `gorm:"type:varchar(255)" json:"photo_url" example:"/uploads/profile/siswa_1.jpg"`
	KelasID      uint             `json:"kelas_id" example:"1"`
	Kelas        kelasModel.Kelas `gorm:"foreignKey:KelasID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"kelas,omitempty"`
	
	// Address Normalization Fields
	Provinsi     string `gorm:"type:varchar(100)" json:"provinsi" example:"Jawa Tengah"`
	Kabupaten    string `gorm:"type:varchar(100)" json:"kabupaten" example:"Salatiga"`
	Kecamatan    string `gorm:"type:varchar(100)" json:"kecamatan" example:"Sidorejo"`
	Desa         string `gorm:"type:varchar(100)" json:"desa" example:"Sidorejo Lor"`
	AlamatDetail string `gorm:"type:text" json:"alamat_detail" example:"Jl. Diponegoro No. 25"`

	// Deprecated Database Columns (Retained for Non-Destructive Migrations)
	TempatLahir string  `gorm:"type:varchar(80)" json:"tempat_lahir,omitempty" example:"Salatiga"`
	Alamat      string  `gorm:"type:text" json:"alamat,omitempty" example:"Jl. Diponegoro No. 25"`
	Email       *string `gorm:"type:varchar(120);uniqueIndex" json:"email,omitempty" example:"budi.santoso@sekolah.com"`

	CreatedAt time.Time      `json:"created_at" example:"2026-07-15T19:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2026-07-15T19:00:00Z"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// SiswaResponseSingle represents a successful single Siswa response
type SiswaResponseSingle struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Berhasil mengambil detail siswa"`
	Data    Siswa  `json:"data"`
}

// SiswaResponseList represents a successful list of Siswa response
type SiswaResponseList struct {
	Success bool    `json:"success" example:"true"`
	Message string  `json:"message" example:"Berhasil mengambil data siswa"`
	Data    []Siswa `json:"data"`
}

// SiswaResponseCreated represents a successful created Siswa response with temporary password
type SiswaResponseCreated struct {
	Success           bool   `json:"success" example:"true"`
	Message           string `json:"message" example:"Siswa berhasil ditambahkan"`
	Data              Siswa  `json:"data"`
	TemporaryPassword string `json:"temporary_password" example:"Siswa123!"`
}
