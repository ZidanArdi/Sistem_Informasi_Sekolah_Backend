package model

import (
	"time"

	authModel "backend/modules/auth/model"

	"gorm.io/gorm"
)

type Guru struct {
	ID           uint            `gorm:"primaryKey" json:"id" example:"1"`
	UserID       uint            `gorm:"index" json:"user_id" example:"3"`
	User         *authModel.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	NIP          string          `gorm:"column:nip;type:varchar(30);uniqueIndex;not null" json:"nip" example:"GR001"`
	Nama         string          `gorm:"type:varchar(100);not null" json:"nama" example:"Ahmad Guru"`
	Gelar        string          `gorm:"type:varchar(50);not null" json:"gelar" example:"S.Pd."`
	JenisKelamin string          `gorm:"type:varchar(20);not null" json:"jenis_kelamin" example:"Laki-laki"`
	NoHP         string          `gorm:"type:varchar(30)" json:"no_hp" example:"081234567890"`
	PhotoURL     string          `gorm:"type:varchar(255)" json:"photo_url" example:"/uploads/profile/guru_1.jpg"`
	
	// Address Normalization Fields
	Provinsi     string `gorm:"type:varchar(100)" json:"provinsi" example:"Jawa Tengah"`
	Kabupaten    string `gorm:"type:varchar(100)" json:"kabupaten" example:"Salatiga"`
	Kecamatan    string `gorm:"type:varchar(100)" json:"kecamatan" example:"Sidorejo"`
	Desa         string `gorm:"type:varchar(100)" json:"desa" example:"Sidorejo Lor"`
	AlamatDetail string `gorm:"type:text" json:"alamat_detail" example:"Jl. Diponegoro No. 25"`

	// Deprecated Database Columns (Retained for Non-Destructive Migrations)
	Alamat string `gorm:"type:text" json:"alamat,omitempty" example:"Jl. Diponegoro No. 25"` // Deprecated

	CreatedAt time.Time      `json:"created_at" example:"2026-07-15T19:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2026-07-15T19:00:00Z"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Virtual field for frontend multiselect checkboxes
	MapelIDs  []uint         `gorm:"-" json:"mapel_ids,omitempty" example:"1,2"`
}

type GuruMapel struct {
	ID      uint `gorm:"primaryKey" json:"id" example:"1"`
	GuruID  uint `gorm:"uniqueIndex:idx_guru_mapel;not null" json:"guru_id" example:"1"`
	MapelID uint `gorm:"uniqueIndex:idx_guru_mapel;not null" json:"mapel_id" example:"2"`
}

func (GuruMapel) TableName() string {
	return "guru_mapel"
}

// GuruResponseSingle represents a successful single Guru response
type GuruResponseSingle struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Berhasil mengambil detail guru"`
	Data    Guru   `json:"data"`
}

// GuruResponseList represents a successful list of Guru response
type GuruResponseList struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Berhasil mengambil data guru"`
	Data    []Guru `json:"data"`
}

// GuruResponseCreated represents a successful created Guru response with temporary password
type GuruResponseCreated struct {
	Success           bool   `json:"success" example:"true"`
	Message           string `json:"message" example:"Guru berhasil ditambahkan"`
	Data              Guru   `json:"data"`
	TemporaryPassword string `json:"temporary_password" example:"Guru123!"`
}
