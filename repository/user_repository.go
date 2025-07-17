package repository

import (
	"hospital-platform/database"
	"hospital-platform/model"
)

// UserRepository kullanıcı veritabanı işlemlerini yönetir
type UserRepository struct{}

// NewUserRepository yeni bir kullanıcı repository'si oluşturur
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create veritabanında yeni bir kullanıcı oluşturur
func (r *UserRepository) Create(user *model.User) error {
	result := database.DB.Create(user)
	return result.Error
}

// GetByID ID'ye göre kullanıcı bilgilerini getirir
func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	result := database.DB.Preload("Hospital").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail e-posta adresine göre kullanıcı olup olmadığını kontrol eder
func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetByPhone telefon numarasına göre kullanıcı olup olmadığını kontrol eder
func (r *UserRepository) GetByPhone(phone string) (*model.User, error) {
	var user model.User
	result := database.DB.Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetByTCKN TC kimlik numarasına göre kullanıcı olup olmadığını kontrol eder
func (r *UserRepository) GetByTCKN(tckn string) (*model.User, error) {
	var user model.User
	result := database.DB.Where("tckn = ?", tckn).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Update kullanıcı bilgilerini günceller
func (r *UserRepository) Update(user *model.User) error {
	result := database.DB.Save(user)
	return result.Error
}

// Delete kullanıcıyı soft delete yapar
func (r *UserRepository) Delete(id uint) error {
	result := database.DB.Delete(&model.User{}, id)
	return result.Error
}

// GetByHospitalID hastaneye ait tüm kullanıcıları getirir
func (r *UserRepository) GetByHospitalID(hospitalID uint) ([]model.User, error) {
	var users []model.User
	result := database.DB.Where("hospital_id = ?", hospitalID).Find(&users)
	return users, result.Error
}
