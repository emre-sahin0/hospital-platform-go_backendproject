package service

import (
	"errors"
	"fmt"
	"hospital-platform/database"
	"hospital-platform/model"
	"hospital-platform/repository"
	"hospital-platform/utils"
)

// Kayıt servisi — user şifresi hashlenir, veritabanına gönderilir
func RegisterUser(user *model.User) error {
	fmt.Println("=== REGISTER DEBUG ===")
	fmt.Println("Kayıt için gelen şifre:", user.Password)

	// Şifreyi güvenli hale getir
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		fmt.Println("Hash hatası:", err)
		return err
	}

	fmt.Println("Hash'lenmiş şifre:", hashedPassword)
	user.Password = hashedPassword

	// Veritabanına kaydet
	err = repository.NewUserRepository().Create(user)
	if err != nil {
		fmt.Println("Veritabanı kayıt hatası:", err)
	} else {
		fmt.Println("Kullanıcı başarıyla kaydedildi")
	}
	fmt.Println("=== REGISTER DEBUG END ===")
	return err
}

func Login(email, password string) (string, error) {
	fmt.Println("=== LOGIN DEBUG ===")
	fmt.Println("Login için gelen email:", email)
	fmt.Println("Login için gelen şifre:", password)

	var user model.User

	// Email ile kullanıcıyı bul
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		fmt.Println("Kullanıcı bulunamadı hatası:", err)
		return "", errors.New("Kullanıcı bulunamadı")
	}

	fmt.Println("Veritabanından gelen hash:", user.Password)
	fmt.Println("Kullanıcının girdiği şifre:", password)

	// Şifreyi kontrol et
	isValid := utils.CheckPasswordHash(password, user.Password)
	fmt.Println("Şifre kontrol sonucu:", isValid)

	if !isValid {
		fmt.Println("Şifre yanlış!")
		return "", errors.New("Şifre yanlış")
	}

	// JWT token üret (doğru parametreler ile)
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role, user.HospitalID, user.Email)
	if err != nil {
		fmt.Println("Token üretme hatası:", err)
		return "", err
	}

	fmt.Println("Login başarılı, token üretildi")
	fmt.Println("=== LOGIN DEBUG END ===")
	return token, nil
}
