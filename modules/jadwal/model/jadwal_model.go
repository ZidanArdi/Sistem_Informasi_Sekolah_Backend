package model

import (
	"time"

	guruModel "backend/modules/guru/model"
	kelasModel "backend/modules/kelas/model"
	mapelModel "backend/modules/mapel/model"

	"gorm.io/gorm"
)

type Jadwal struct {
	ID         uint             `gorm:"primaryKey" json:"id"`
	KelasID    uint             `gorm:"not null" json:"kelas_id"`
	MapelID    uint             `gorm:"not null" json:"mapel_id"`
	GuruID     uint             `gorm:"not null" json:"guru_id"`
	Hari       string           `gorm:"type:varchar(20);not null" json:"hari"`
	JamMulai   string           `gorm:"type:time;not null" json:"jam_mulai"`
	JamSelesai string           `gorm:"type:time;not null" json:"jam_selesai"`
	Kelas      kelasModel.Kelas `gorm:"foreignKey:KelasID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"kelas,omitempty"`
	Mapel      mapelModel.Mapel `gorm:"foreignKey:MapelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"mapel,omitempty"`
	Guru       guruModel.Guru   `gorm:"foreignKey:GuruID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"guru,omitempty"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	DeletedAt  gorm.DeletedAt   `gorm:"index" json:"-"`
}
