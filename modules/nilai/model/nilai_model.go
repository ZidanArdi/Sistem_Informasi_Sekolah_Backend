package model

import (
	"time"

	mapelModel "backend/modules/mapel/model"
	siswaModel "backend/modules/siswa/model"

	"gorm.io/gorm"
)

type Nilai struct {
	ID         uint             `gorm:"primaryKey" json:"id"`
	SiswaID    uint             `gorm:"not null" json:"siswa_id"`
	MapelID    uint             `gorm:"not null" json:"mapel_id"`
	Semester   string           `gorm:"type:varchar(20);not null" json:"semester"`
	JenisNilai string           `gorm:"type:varchar(20);not null" json:"jenis_nilai"`
	Nilai      float64          `gorm:"not null" json:"nilai"`
	Siswa      siswaModel.Siswa `gorm:"foreignKey:SiswaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"siswa,omitempty"`
	Mapel      mapelModel.Mapel `gorm:"foreignKey:MapelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"mapel,omitempty"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	DeletedAt  gorm.DeletedAt   `gorm:"index" json:"-"`
}
