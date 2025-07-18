package service

import (
	"fmt"
	"hospital-platform/database"
	"hospital-platform/model"
	"hospital-platform/repository"
	"hospital-platform/utils"
)

// HospitalService hastane iş mantığını yönetir
type HospitalService struct {
	hospitalRepo *repository.HospitalRepository
	locationRepo *repository.LocationRepository
	userRepo     *repository.UserRepository
}

// NewHospitalService yeni bir hastane servisi oluşturur
func NewHospitalService() *HospitalService {
	return &HospitalService{
		hospitalRepo: repository.NewHospitalRepository(),
		locationRepo: repository.NewLocationRepository(),
		userRepo:     repository.NewUserRepository(),
	}
}

// RegisterHospital yeni hastane ve admin kullanıcı kaydı yapar
func (s *HospitalService) RegisterHospital(req *model.HospitalRegistrationRequest) (*model.HospitalRegistrationResponse, []model.ValidationError, error) {
	// 1. Giriş verilerini doğrula
	validationErrors := s.validateRegistrationData(req)
	if len(validationErrors) > 0 {
		return nil, validationErrors, nil
	}

	// 2. Veritabanı transaction'ını başlat
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, nil, fmt.Errorf("transaction başlatılamadı: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 3. Hastaneyi oluştur
	hospital := &model.Hospital{
		Name:          req.HospitalName,
		TaxID:         req.TaxID,
		Email:         req.HospitalEmail,
		Phone:         req.HospitalPhone,
		ProvinceID:    req.ProvinceID,
		DistrictID:    req.DistrictID,
		AddressDetail: req.AddressDetail,
	}

	if err := tx.Create(hospital).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("hastane oluşturulamadı: %v", err)
	}

	// 4. Admin şifresini hash'le
	hashedPassword, err := utils.HashPassword(req.AdminPassword)
	if err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("şifre hash'lenemedi: %v", err)
	}

	// 5. Admin kullanıcıyı oluştur
	adminUser := &model.User{
		HospitalID: hospital.ID,
		FirstName:  req.AdminFirstName,
		LastName:   req.AdminLastName,
		TCKN:       req.AdminTCKN,
		Email:      req.AdminEmail,
		Phone:      req.AdminPhone,
		Password:   hashedPassword,
		Role:       model.RoleYetkili, // İlk kullanıcı yetkili
		CreatedBy:  nil,               // İlk kullanıcı için nil
		IsActive:   true,
	}

	if err := tx.Create(adminUser).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("admin kullanıcı oluşturulamadı: %v", err)
	}

	// 6. Transaction'ı commit et
	if err := tx.Commit().Error; err != nil {
		return nil, nil, fmt.Errorf("transaction commit edilemedi: %v", err)
	}

	// 7. JWT token oluştur (hospital_id ve username dahil)
	// Email'i username olarak kullan
	token, err := utils.GenerateJWT(adminUser.ID, adminUser.Email, adminUser.Role, adminUser.HospitalID, adminUser.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("token oluşturulamadı: %v", err)
	}

	// 8. Response için ilişkileri yükle
	hospital, _ = s.hospitalRepo.GetByID(hospital.ID)
	adminUser.Password = "" // Şifreyi response'dan kaldır

	response := &model.HospitalRegistrationResponse{
		Message:   "Hastane ve admin kullanıcı başarıyla oluşturuldu",
		Hospital:  *hospital,
		AdminUser: *adminUser,
		Token:     token,
	}

	return response, nil, nil
}

// validateRegistrationData hastane kayıt verilerini doğrular
func (s *HospitalService) validateRegistrationData(req *model.HospitalRegistrationRequest) []model.ValidationError {
	var errors []model.ValidationError

	// Hastane vergi kimlik numarası benzersizlik kontrolü
	if existingHospital, _ := s.hospitalRepo.GetByTaxID(req.TaxID); existingHospital != nil {
		errors = append(errors, model.ValidationError{
			Field:   "tax_id",
			Message: "Bu vergi kimlik numarası zaten kullanılıyor",
		})
	}

	// Hastane e-posta benzersizlik kontrolü
	if existingHospital, _ := s.hospitalRepo.GetByEmail(req.HospitalEmail); existingHospital != nil {
		errors = append(errors, model.ValidationError{
			Field:   "hospital_email",
			Message: "Bu e-posta adresi zaten kullanılıyor",
		})
	}

	// Hastane telefon benzersizlik kontrolü
	if existingHospital, _ := s.hospitalRepo.GetByPhone(req.HospitalPhone); existingHospital != nil {
		errors = append(errors, model.ValidationError{
			Field:   "hospital_phone",
			Message: "Bu telefon numarası zaten kullanılıyor",
		})
	}

	// İl-ilçe uyumluluk kontrolü
	isValid, err := s.locationRepo.ValidateProvinceDistrict(req.ProvinceID, req.DistrictID)
	if err != nil || !isValid {
		errors = append(errors, model.ValidationError{
			Field:   "district_id",
			Message: "Seçilen ilçe, seçilen ile ait değil",
		})
	}

	// Admin TC kimlik numarası benzersizlik kontrolü
	if existingUser, _ := s.userRepo.GetByTCKN(req.AdminTCKN); existingUser != nil {
		errors = append(errors, model.ValidationError{
			Field:   "admin_tc",
			Message: "Bu TC kimlik numarası zaten kullanılıyor",
		})
	}

	// Admin e-posta benzersizlik kontrolü
	if existingUser, _ := s.userRepo.GetByEmail(req.AdminEmail); existingUser != nil {
		errors = append(errors, model.ValidationError{
			Field:   "admin_email",
			Message: "Bu e-posta adresi zaten kullanılıyor",
		})
	}

	// Admin telefon benzersizlik kontrolü
	if existingUser, _ := s.userRepo.GetByPhone(req.AdminPhone); existingUser != nil {
		errors = append(errors, model.ValidationError{
			Field:   "admin_phone",
			Message: "Bu telefon numarası zaten kullanılıyor",
		})
	}

	return errors
}

// GetHospitalByID ID'ye göre hastane bilgilerini getirir
func (s *HospitalService) GetHospitalByID(id uint) (*model.Hospital, error) {
	return s.hospitalRepo.GetByID(id)
}
