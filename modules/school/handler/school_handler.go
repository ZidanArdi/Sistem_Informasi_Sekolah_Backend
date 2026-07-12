package handler

import (
	"backend/config"
	"backend/helpers"
	"backend/modules/school/model"

	"github.com/gofiber/fiber/v2"
)

func GetSchoolProfile(c *fiber.Ctx) error {
	var profile model.SchoolProfile
	// Always look up the record with ID = 1
	err := config.DB.First(&profile, 1).Error
	if err != nil {
		// If not found in database, return a default profile struct
		defaultProfile := model.SchoolProfile{
			ID:                1,
			Logo:              "",
			Name:              "SMK Negeri 1 Salatiga",
			Npsn:              "20312345",
			Status:            "Negeri",
			Level:             "SMK / Sekolah Menengah Kejuruan",
			Accreditation:     "A (Sangat Baik)",
			EstablishedYear:   "1965",
			Address:           "Jl. Diponegoro No. 25, Salatiga",
			PostalCode:        "50711",
			Phone:             "(0298) 123456",
			Email:             "info@smkn1salatiga.sch.id",
			Website:           "https://smkn1salatiga.sch.id",
			PrincipalName:     "Drs. Hadi Santoso",
			PrincipalNip:      "197205121998031002",
			PrincipalPosition: "Kepala Sekolah",
			AppointmentPeriod: "2021 - 2027",
			AcademicYear:      "2026/2027",
			CurrentSemester:   "Ganjil",
			AcademicStatus:    "🟢 Aktif",
			SchoolType:        "Sekolah Menengah Kejuruan (SMK)",
			Curriculum:        "Kurikulum Merdeka",
			Shift:             "Pagi",
			OperationalStatus: "Aktif",
			Vision:            "Menjadi lembaga pendidikan kejuruan yang unggul, berkarakter, dan berdaya saing global di era digital.",
			Mission:           "Menyelenggarakan pembelajaran berbasis kompetensi teknologi industri dan keahlian terkini.\nMembina karakter siswa berlandaskan iman, taqwa, nilai Pancasila, dan budi pekerti luhur.\nMenjalin kemitraan erat dengan dunia usaha, dunia industri, dan asosiasi profesi (DUDI).\nMengembangkan semangat kewirausahaan, daya saing, dan kemandirian siswa.",
			Facilities:        "Library, Computer Laboratory, Mosque, Workshop, Sports Field",
		}
		return helpers.SuccessResponse(c, "Berhasil mengambil profil default", defaultProfile)
	}

	return helpers.SuccessResponse(c, "Berhasil mengambil profil sekolah", profile)
}

func UpdateSchoolProfile(c *fiber.Ctx) error {
	var input model.SchoolProfile
	if err := c.BodyParser(&input); err != nil {
		return helpers.ErrorResponse(c, 400, "Input tidak valid")
	}

	var profile model.SchoolProfile
	err := config.DB.First(&profile, 1).Error
	if err != nil {
		// If not found, create new with ID 1
		input.ID = 1
		if errCreate := config.DB.Create(&input).Error; errCreate != nil {
			return helpers.ErrorResponse(c, 500, "Gagal membuat profil sekolah: "+errCreate.Error())
		}
		return helpers.SuccessResponse(c, "Profil sekolah berhasil dibuat", input)
	}

	// Update existing record
	profile.Logo = input.Logo
	profile.Name = input.Name
	profile.Npsn = input.Npsn
	profile.Status = input.Status
	profile.Level = input.Level
	profile.Accreditation = input.Accreditation
	profile.EstablishedYear = input.EstablishedYear
	profile.Address = input.Address
	profile.PostalCode = input.PostalCode
	profile.Phone = input.Phone
	profile.Email = input.Email
	profile.Website = input.Website
	profile.PrincipalName = input.PrincipalName
	profile.PrincipalNip = input.PrincipalNip
	profile.PrincipalPosition = input.PrincipalPosition
	profile.AppointmentPeriod = input.AppointmentPeriod
	profile.AcademicYear = input.AcademicYear
	profile.CurrentSemester = input.CurrentSemester
	profile.AcademicStatus = input.AcademicStatus
	profile.SchoolType = input.SchoolType
	profile.Curriculum = input.Curriculum
	profile.Shift = input.Shift
	profile.OperationalStatus = input.OperationalStatus
	profile.Vision = input.Vision
	profile.Mission = input.Mission
	profile.Facilities = input.Facilities

	if errSave := config.DB.Save(&profile).Error; errSave != nil {
		return helpers.ErrorResponse(c, 500, "Gagal memperbarui profil sekolah: "+errSave.Error())
	}

	return helpers.SuccessResponse(c, "Profil sekolah berhasil diperbarui", profile)
}
