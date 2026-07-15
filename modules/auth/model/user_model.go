package model

import "time"

type User struct {
	ID           uint       `gorm:"primaryKey" json:"id" example:"1"`
	Nama         string     `gorm:"type:varchar(100);not null" json:"nama" example:"Ahmad Guru"`
	Email        *string    `gorm:"type:varchar(120);uniqueIndex" json:"email" example:"ahmad.guru@sekolah.com"`
	Password     string     `gorm:"type:varchar(255);not null" json:"-"`
	Role         string     `gorm:"type:varchar(20);not null;default:siswa" json:"role" example:"guru"`
	IsFirstLogin bool       `gorm:"type:boolean;not null;default:false" json:"is_first_login" example:"true"`
	IsActive     bool       `gorm:"type:boolean;not null;default:true" json:"is_active" example:"true"`
	LastLoginAt  *time.Time `gorm:"type:timestamp" json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at" example:"2026-07-15T19:00:00Z"`
	UpdatedAt    time.Time  `json:"updated_at" example:"2026-07-15T19:00:00Z"`

	// Ignored by GORM, populated for unified API login response
	NIP          string `json:"nip,omitempty" gorm:"-" example:"GR001"`
	NIS          string `json:"nis,omitempty" gorm:"-" example:"12345678"`
	NoHP         string `json:"no_hp,omitempty" gorm:"-" example:"081234567890"`
	Alamat       string `json:"alamat,omitempty" gorm:"-" example:"Jl. Diponegoro No. 25, Salatiga"`
	Gelar        string `json:"gelar,omitempty" gorm:"-" example:"S.Pd."`
	JenisKelamin string `json:"jenis_kelamin,omitempty" gorm:"-" example:"Laki-laki"`
}

func (User) TableName() string {
	return "users"
}

// RegisterSuccessResponse represents a successful 201 Created register response
type RegisterSuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Register berhasil"`
	Data    User   `json:"data"`
}

// LoginData represents the payload of a successful login response
type LoginData struct {
	User  User   `json:"user"`
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6ImFobWFkLmd1cnVAc2Vrb2xhaC5jb20iLCJyb2xlIjoiZ3VydSIsImV4cCI6MTcyMTE0MzQwMH0.xxxx"`
}

// LoginSuccessResponse represents a successful 200 OK login response
type LoginSuccessResponse struct {
	Success bool      `json:"success" example:"true"`
	Message string    `json:"message" example:"Login berhasil"`
	Data    LoginData `json:"data"`
}
