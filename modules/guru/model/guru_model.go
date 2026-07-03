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
	Alamat       string          `gorm:"type:text" json:"alamat"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    gorm.DeletedAt  `gorm:"index" json:"-"`
}
