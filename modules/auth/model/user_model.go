package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nama      string    `gorm:"type:varchar(100);not null" json:"nama"`
	Email        *string   `gorm:"type:varchar(120);uniqueIndex" json:"email"`
	Password     string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         string    `gorm:"type:varchar(20);not null;default:siswa" json:"role"`
	IsFirstLogin bool      `gorm:"type:boolean;not null;default:false" json:"is_first_login"`
	IsActive     bool      `gorm:"type:boolean;not null;default:true" json:"is_active"`
	LastLoginAt  *time.Time `gorm:"type:timestamp" json:"last_login_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Ignored by GORM, populated for unified API login response
	NIP          string `json:"nip,omitempty" gorm:"-"`
	NIS          string `json:"nis,omitempty" gorm:"-"`
	NoHP         string `json:"no_hp,omitempty" gorm:"-"`
	Alamat       string `json:"alamat,omitempty" gorm:"-"`
	Gelar        string `json:"gelar,omitempty" gorm:"-"`
	JenisKelamin string `json:"jenis_kelamin,omitempty" gorm:"-"`
}

func (User) TableName() string {
	return "users"
}
