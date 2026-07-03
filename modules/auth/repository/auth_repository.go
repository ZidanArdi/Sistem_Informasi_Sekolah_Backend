package repository

import (
	"backend/config"
	"backend/modules/auth/model"

	"gorm.io/gorm"
)

func CreateUser(user model.User) (model.User, error) {
	result := config.DB.Create(&user)
	return user, result.Error
}

func GetUserByEmail(email string) (model.User, error) {
	var user model.User
	result := config.DB.Where("email = ?", email).First(&user)
	return user, result.Error
}

func GetUserByID(id uint) (model.User, error) {
	var user model.User
	result := config.DB.First(&user, id)
	return user, result.Error
}

func UpdatePassword(id uint, password string) error {
	return config.DB.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"password":       password,
		"is_first_login": false,
	}).Error
}

func GetUserIDByNIS(nis string) (uint, string, error) {
	var result struct {
		UserID uint   `gorm:"column:user_id"`
		Email  string `gorm:"column:email"`
	}
	err := config.DB.Table("siswas").Select("user_id, email").Where("nis = ? AND deleted_at IS NULL", nis).First(&result).Error
	return result.UserID, result.Email, err
}

func EmailExists(email string) bool {
	var count int64
	config.DB.Model(&model.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func IsNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}
