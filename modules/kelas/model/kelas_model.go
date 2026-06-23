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
	WaliKelasID *uint           `gorm:"uniqueIndex" json:"wali_kelas_id"`
	WaliKelas   *guruModel.Guru `gorm:"foreignKey:WaliKelasID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"wali_kelas,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}
