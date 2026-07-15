package model

import (
	"time"

	"gorm.io/gorm"
)

type SchoolProfile struct {
	ID                uint           `gorm:"primaryKey" json:"id" example:"1"`
	Logo              string         `gorm:"type:varchar(255)" json:"logo" example:"/uploads/profile/logo.png"`
	Name              string         `gorm:"type:varchar(150);not null" json:"name" example:"SMK Negeri 1 Salatiga"`
	Npsn              string         `gorm:"type:varchar(30);not null" json:"npsn" example:"20312345"`
	Status            string         `gorm:"type:varchar(30);not null" json:"status" example:"Negeri"`
	Level             string         `gorm:"type:varchar(100);not null" json:"level" example:"SMK / Sekolah Menengah Kejuruan"`
	Accreditation     string         `gorm:"type:varchar(50);not null" json:"accreditation" example:"A (Sangat Baik)"`
	EstablishedYear   string         `gorm:"type:varchar(10);not null" json:"established_year" example:"1965"`
	Address           string         `gorm:"type:text;not null" json:"address" example:"Jl. Diponegoro No. 25, Salatiga"`
	PostalCode        string         `gorm:"type:varchar(20);not null" json:"postal_code" example:"50711"`
	Phone             string         `gorm:"type:varchar(30);not null" json:"phone" example:"(0298) 123456"`
	Email             string         `gorm:"type:varchar(100);not null" json:"email" example:"info@smkn1salatiga.sch.id"`
	Website           string         `gorm:"type:varchar(150);not null" json:"website" example:"https://smkn1salatiga.sch.id"`
	PrincipalName     string         `gorm:"type:varchar(100);not null" json:"principal_name" example:"Drs. Hadi Santoso"`
	PrincipalNip      string         `gorm:"type:varchar(30);not null" json:"principal_nip" example:"197205121998031002"`
	PrincipalPosition string         `gorm:"type:varchar(100);not null" json:"principal_position" example:"Kepala Sekolah"`
	AppointmentPeriod string         `gorm:"type:varchar(50);not null" json:"appointment_period" example:"2021 - 2027"`
	AcademicYear      string         `gorm:"type:varchar(30);not null" json:"academic_year" example:"2026/2027"`
	CurrentSemester   string         `gorm:"type:varchar(20);not null" json:"current_semester" example:"Ganjil"`
	AcademicStatus    string         `gorm:"type:varchar(30);not null" json:"academic_status" example:"🟢 Aktif"`
	SchoolType        string         `gorm:"type:varchar(100);not null" json:"school_type" example:"Sekolah Menengah Kejuruan (SMK)"`
	Curriculum        string         `gorm:"type:varchar(100);not null" json:"curriculum" example:"Kurikulum Merdeka"`
	Shift             string         `gorm:"type:varchar(30);not null" json:"shift" example:"Pagi"`
	OperationalStatus string         `gorm:"type:varchar(30);not null" json:"operational_status" example:"Aktif"`
	Vision            string         `gorm:"type:text;not null" json:"vision" example:"Menjadi lembaga pendidikan kejuruan yang unggul, berkarakter, dan berdaya saing global."`
	Mission           string         `gorm:"type:text;not null" json:"mission" example:"1. Menyelenggarakan pembelajaran berkualitas...\n2. Membina budi pekerti luhur..."`
	Facilities        string         `gorm:"type:text;not null" json:"facilities" example:"Perpustakaan, Lab Komputer, Bengkel Kerja, Lapangan Olahraga"`
	CreatedAt         time.Time      `json:"created_at" example:"2026-07-15T19:00:00Z"`
	UpdatedAt         time.Time      `json:"updated_at" example:"2026-07-15T19:00:00Z"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// SchoolResponseSingle represents a successful SchoolProfile response
type SchoolResponseSingle struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message" example:"Berhasil mengambil profil sekolah"`
	Data    SchoolProfile `json:"data"`
}
