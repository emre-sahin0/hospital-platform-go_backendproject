package repository

import (
	"errors"
	"hospital-platform/database"
	"hospital-platform/model"

	"gorm.io/gorm"
)

// Veritabanına kullanıcıyı kaydeder
func CreateUser(user *model.User) error {
	result := database.DB.Create(user)
	return result.Error
}

func GetUserByPhone(phone string) (*model.User, error) {
	var user model.User
	if err := database.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // kullanıcı bulunamadı
		}
		return nil, err // başka hata
	}
	return &user, nil
}
