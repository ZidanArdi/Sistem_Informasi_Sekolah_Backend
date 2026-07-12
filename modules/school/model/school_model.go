package model

import (
	"time"

	"gorm.io/gorm"
)

type SchoolProfile struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	Logo              string         `gorm:"type:varchar(255)" json:"logo"`
	Name              string         `gorm:"type:varchar(150);not null" json:"name"`
	Npsn              string         `gorm:"type:varchar(30);not null" json:"npsn"`
	Status            string         `gorm:"type:varchar(30);not null" json:"status"`
	Level             string         `gorm:"type:varchar(100);not null" json:"level"`
	Accreditation     string         `gorm:"type:varchar(50);not null" json:"accreditation"`
	EstablishedYear   string         `gorm:"type:varchar(10);not null" json:"established_year"`
	Address           string         `gorm:"type:text;not null" json:"address"`
	PostalCode        string         `gorm:"type:varchar(20);not null" json:"postal_code"`
	Phone             string         `gorm:"type:varchar(30);not null" json:"phone"`
	Email             string         `gorm:"type:varchar(100);not null" json:"email"`
	Website           string         `gorm:"type:varchar(150);not null" json:"website"`
	PrincipalName     string         `gorm:"type:varchar(100);not null" json:"principal_name"`
	PrincipalNip      string         `gorm:"type:varchar(30);not null" json:"principal_nip"`
	PrincipalPosition string         `gorm:"type:varchar(100);not null" json:"principal_position"`
	AppointmentPeriod string         `gorm:"type:varchar(50);not null" json:"appointment_period"`
	AcademicYear      string         `gorm:"type:varchar(30);not null" json:"academic_year"`
	CurrentSemester   string         `gorm:"type:varchar(20);not null" json:"current_semester"`
	AcademicStatus    string         `gorm:"type:varchar(30);not null" json:"academic_status"`
	SchoolType        string         `gorm:"type:varchar(100);not null" json:"school_type"`
	Curriculum        string         `gorm:"type:varchar(100);not null" json:"curriculum"`
	Shift             string         `gorm:"type:varchar(30);not null" json:"shift"`
	OperationalStatus string         `gorm:"type:varchar(30);not null" json:"operational_status"`
	Vision            string         `gorm:"type:text;not null" json:"vision"`
	Mission           string         `gorm:"type:text;not null" json:"mission"`
	Facilities        string         `gorm:"type:text;not null" json:"facilities"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}
