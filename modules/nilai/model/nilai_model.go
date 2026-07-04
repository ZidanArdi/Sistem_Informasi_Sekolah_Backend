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
	ID          uint             `gorm:"primaryKey" json:"id"`
	SiswaID     uint             `gorm:"uniqueIndex:idx_nilai_unique,priority:1;not null" json:"siswa_id"`
	KelasID     uint             `gorm:"not null" json:"kelas_id"`
	MapelID     uint             `gorm:"uniqueIndex:idx_nilai_unique,priority:2;not null" json:"mapel_id"`
	GuruID      uint             `gorm:"not null" json:"guru_id"`
	Tugas       float64          `gorm:"type:decimal(5,2);not null;default:0" json:"tugas"`
	UTS         float64          `gorm:"type:decimal(5,2);not null;default:0" json:"uts"`
	UAS         float64          `gorm:"type:decimal(5,2);not null;default:0" json:"uas"`
	NilaiAkhir  float64          `gorm:"type:decimal(5,2);not null;default:0" json:"nilai_akhir"`
	GradeHuruf  string           `gorm:"type:varchar(5);not null" json:"grade_huruf"`
	Semester    string           `gorm:"type:varchar(20);uniqueIndex:idx_nilai_unique,priority:3;not null" json:"semester"`
	TahunAjaran string           `gorm:"type:varchar(20);uniqueIndex:idx_nilai_unique,priority:4;not null" json:"tahun_ajaran"`
	Siswa       siswaModel.Siswa `gorm:"foreignKey:SiswaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"siswa,omitempty"`
	Kelas       kelasModel.Kelas `gorm:"foreignKey:KelasID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"kelas,omitempty"`
	Mapel       mapelModel.Mapel `gorm:"foreignKey:MapelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"mapel,omitempty"`
	Guru        guruModel.Guru   `gorm:"foreignKey:GuruID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"guru,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `gorm:"index" json:"-"`
}
