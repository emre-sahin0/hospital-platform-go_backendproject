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

func Login(emailOrPhone, password string) (string, error) {
	fmt.Println("=== LOGIN DEBUG ===")
	fmt.Println("Login için gelen email/telefon:", emailOrPhone)
	fmt.Println("Login için gelen şifre:", password)

	var user model.User

	// Email veya telefon ile kullanıcıyı bul
	if err := database.DB.Where("email = ? OR phone = ?", emailOrPhone, emailOrPhone).First(&user).Error; err != nil {
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

// ==================== ALT KULLANICI YÖNETİMİ ====================

// CreateSubUser alt kullanıcı oluşturur
func CreateSubUser(req *model.CreateSubUserRequest, createdBy uint) (*model.User, []model.ValidationError, error) {
	// 1. Giriş verilerini doğrula
	validationErrors := validateSubUserData(req)
	if len(validationErrors) > 0 {
		return nil, validationErrors, nil
	}

	// 2. Oluşturan kullanıcının hastane ID'sini al
	var creator model.User
	if err := database.DB.First(&creator, createdBy).Error; err != nil {
		return nil, nil, fmt.Errorf("oluşturan kullanıcı bulunamadı: %v", err)
	}

	// 3. Şifreyi hash'le
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, nil, fmt.Errorf("şifre hash'lenemedi: %v", err)
	}

	// 4. Alt kullanıcıyı oluştur
	subUser := &model.User{
		HospitalID: creator.HospitalID,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		TCKN:       req.TCKN,
		Email:      req.Email,
		Phone:      req.Phone,
		Password:   hashedPassword,
		Role:       req.Role,
		CreatedBy:  &createdBy,
		IsActive:   true,
	}

	if err := database.DB.Create(subUser).Error; err != nil {
		return nil, nil, fmt.Errorf("alt kullanıcı oluşturulamadı: %v", err)
	}

	return subUser, nil, nil
}

// GetSubUsers hastaneye ait alt kullanıcıları getirir
func GetSubUsers(hospitalID uint) ([]model.User, error) {
	var users []model.User
	if err := database.DB.Where("hospital_id = ?", hospitalID).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("kullanıcılar getirilemedi: %v", err)
	}
	return users, nil
}

// UpdateSubUser alt kullanıcı bilgilerini günceller
func UpdateSubUser(userID uint, req *model.UpdateSubUserRequest, updatedBy uint) (*model.User, []model.ValidationError, error) {
	// 1. Giriş verilerini doğrula
	validationErrors := validateUpdateSubUserData(req, userID)
	if len(validationErrors) > 0 {
		return nil, validationErrors, nil
	}

	// 2. Kullanıcıyı bul
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, nil, fmt.Errorf("kullanıcı bulunamadı: %v", err)
	}

	// 3. Güncelleme yetkisi kontrolü
	var updater model.User
	if err := database.DB.First(&updater, updatedBy).Error; err != nil {
		return nil, nil, fmt.Errorf("güncelleyen kullanıcı bulunamadı: %v", err)
	}

	if updater.HospitalID != user.HospitalID {
		return nil, nil, fmt.Errorf("farklı hastaneye ait kullanıcı güncellenemez")
	}

	// 4. Kullanıcı bilgilerini güncelle
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.Phone = req.Phone
	user.Role = req.Role
	user.IsActive = req.IsActive

	if err := database.DB.Save(&user).Error; err != nil {
		return nil, nil, fmt.Errorf("kullanıcı güncellenemedi: %v", err)
	}

	return &user, nil, nil
}

// DeleteSubUser alt kullanıcıyı siler
func DeleteSubUser(userID uint, deletedBy uint) error {
	// 1. Kullanıcıyı bul
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return fmt.Errorf("kullanıcı bulunamadı: %v", err)
	}

	// 2. Silme yetkisi kontrolü
	var deleter model.User
	if err := database.DB.First(&deleter, deletedBy).Error; err != nil {
		return fmt.Errorf("silen kullanıcı bulunamadı: %v", err)
	}

	if deleter.HospitalID != user.HospitalID {
		return fmt.Errorf("farklı hastaneye ait kullanıcı silinemez")
	}

	// 3. Kendini silmeye çalışıyorsa engelle
	if userID == deletedBy {
		return fmt.Errorf("kendinizi silemezsiniz")
	}

	// 4. Kullanıcıyı soft delete yap
	if err := database.DB.Delete(&user).Error; err != nil {
		return fmt.Errorf("kullanıcı silinemedi: %v", err)
	}

	return nil
}

// ==================== VALIDATION FUNCTIONS ====================

// validateSubUserData alt kullanıcı verilerini doğrular
func validateSubUserData(req *model.CreateSubUserRequest) []model.ValidationError {
	var errors []model.ValidationError

	// TC kimlik kontrolü
	if len(req.TCKN) != 11 {
		errors = append(errors, model.ValidationError{
			Field:   "tc",
			Message: "TC kimlik numarası 11 haneli olmalıdır",
		})
	}

	// TC kimlik benzersizlik kontrolü
	var existingUser model.User
	if err := database.DB.Where("tckn = ?", req.TCKN).First(&existingUser).Error; err == nil {
		errors = append(errors, model.ValidationError{
			Field:   "tc",
			Message: "Bu TC kimlik numarası zaten kullanılıyor",
		})
	}

	// Email benzersizlik kontrolü
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		errors = append(errors, model.ValidationError{
			Field:   "email",
			Message: "Bu e-posta adresi zaten kullanılıyor",
		})
	}

	// Telefon benzersizlik kontrolü
	if err := database.DB.Where("phone = ?", req.Phone).First(&existingUser).Error; err == nil {
		errors = append(errors, model.ValidationError{
			Field:   "phone",
			Message: "Bu telefon numarası zaten kullanılıyor",
		})
	}

	return errors
}

// validateUpdateSubUserData güncelleme verilerini doğrular
func validateUpdateSubUserData(req *model.UpdateSubUserRequest, userID uint) []model.ValidationError {
	var errors []model.ValidationError

	// Email benzersizlik kontrolü (kendisi hariç)
	var existingUser model.User
	if err := database.DB.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser).Error; err == nil {
		errors = append(errors, model.ValidationError{
			Field:   "email",
			Message: "Bu e-posta adresi zaten kullanılıyor",
		})
	}

	// Telefon benzersizlik kontrolü (kendisi hariç)
	if err := database.DB.Where("phone = ? AND id != ?", req.Phone, userID).First(&existingUser).Error; err == nil {
		errors = append(errors, model.ValidationError{
			Field:   "phone",
			Message: "Bu telefon numarası zaten kullanılıyor",
		})
	}

	return errors
}
