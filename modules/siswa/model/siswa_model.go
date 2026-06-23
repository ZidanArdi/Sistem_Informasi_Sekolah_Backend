package model

import (
	"time"

	kelasModel "backend/modules/kelas/model"

	"gorm.io/gorm"
)

type Siswa struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	NIS          string           `gorm:"type:varchar(30);uniqueIndex;not null" json:"nis"`
	Nama         string           `gorm:"type:varchar(100);not null" json:"nama"`
	JenisKelamin string           `gorm:"type:varchar(20)" json:"jenis_kelamin"`
	TempatLahir  string           `gorm:"type:varchar(80)" json:"tempat_lahir"`
	TanggalLahir string           `gorm:"type:date" json:"tanggal_lahir"`
	Alamat       string           `gorm:"type:text;not null" json:"alamat"`
	KelasID      uint             `json:"kelas_id"`
	Kelas        kelasModel.Kelas `gorm:"foreignKey:KelasID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"kelas,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"-"`
}
