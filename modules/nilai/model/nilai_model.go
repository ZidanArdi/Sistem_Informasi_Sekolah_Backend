package model

import (
	"time"

	guruModel "backend/modules/guru/model"
	kelasModel "backend/modules/kelas/model"
	mapelModel "backend/modules/mapel/model"
	siswaModel "backend/modules/siswa/model"

	"gorm.io/gorm"
)

type Nilai struct {
	ID          uint             `gorm:"primaryKey" json:"id" example:"1"`
	SiswaID     uint             `gorm:"uniqueIndex:idx_nilai_unique,priority:1;not null" json:"siswa_id" example:"1"`
	KelasID     uint             `gorm:"not null" json:"kelas_id" example:"1"`
	MapelID     uint             `gorm:"uniqueIndex:idx_nilai_unique,priority:2;not null" json:"mapel_id" example:"1"`
	GuruID      uint             `gorm:"not null" json:"guru_id" example:"1"`
	Tugas       float64          `gorm:"type:decimal(5,2);not null;default:0" json:"tugas" example:"80.0"`
	UTS         float64          `gorm:"type:decimal(5,2);not null;default:0" json:"uts" example:"85.0"`
	UAS         float64          `gorm:"type:decimal(5,2);not null;default:0" json:"uas" example:"90.0"`
	NilaiAkhir  float64          `gorm:"type:decimal(5,2);not null;default:0" json:"nilai_akhir" example:"85.5"`
	GradeHuruf  string           `gorm:"type:varchar(5);not null" json:"grade_huruf" example:"A"`
	Semester    string           `gorm:"type:varchar(20);uniqueIndex:idx_nilai_unique,priority:3;not null" json:"semester" example:"Ganjil"`
	TahunAjaran string           `gorm:"type:varchar(20);uniqueIndex:idx_nilai_unique,priority:4;not null" json:"tahun_ajaran" example:"2026/2027"`
	Siswa       siswaModel.Siswa `gorm:"foreignKey:SiswaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"siswa,omitempty"`
	Kelas       kelasModel.Kelas `gorm:"foreignKey:KelasID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"kelas,omitempty"`
	Mapel       mapelModel.Mapel `gorm:"foreignKey:MapelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"mapel,omitempty"`
	Guru        guruModel.Guru   `gorm:"foreignKey:GuruID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"guru,omitempty"`
	CreatedAt   time.Time        `json:"created_at" example:"2026-07-15T19:00:00Z"`
	UpdatedAt   time.Time        `json:"updated_at" example:"2026-07-15T19:00:00Z"`
	DeletedAt   gorm.DeletedAt   `gorm:"index" json:"-"`
}

// NilaiResponseSingle represents a successful single Nilai response
type NilaiResponseSingle struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Berhasil mengambil detail nilai siswa"`
	Data    Nilai  `json:"data"`
}

// NilaiResponseList represents a successful list of Nilai response
type NilaiResponseList struct {
	Success bool    `json:"success" example:"true"`
	Message string  `json:"message" example:"Berhasil mengambil data nilai siswa"`
	Data    []Nilai `json:"data"`
}
