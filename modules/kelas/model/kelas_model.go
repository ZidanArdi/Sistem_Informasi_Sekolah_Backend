package model

import (
	"time"

	guruModel "backend/modules/guru/model"

	"gorm.io/gorm"
)

type Kelas struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	NamaKelas   string          `gorm:"type:varchar(80);not null" json:"nama_kelas"`
	Tingkat     string          `gorm:"type:varchar(20);not null" json:"tingkat"`
	Jurusan     string          `gorm:"type:varchar(50)" json:"jurusan"`
	Kapasitas   int             `gorm:"not null;default:0" json:"kapasitas"`
	WaliKelasID *uint           `gorm:"uniqueIndex" json:"wali_kelas_id"`
	WaliKelas   *guruModel.Guru `gorm:"foreignKey:WaliKelasID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"wali_kelas,omitempty"`
	TotalSiswa  int             `gorm:"-" json:"total_siswa"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}
