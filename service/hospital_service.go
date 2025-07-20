package service

import (
	"fmt"
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

// RegisterHospital hastane ve admin kullanıcı kaydı yapar
func (s *HospitalService) RegisterHospital(req *model.HospitalRegistrationRequest) (*model.HospitalRegistrationResponse, []model.ValidationError, error) {
	// 1. Giriş verilerini doğrula
	validationErrors := s.validateRegistrationData(req)
	if len(validationErrors) > 0 {
		return nil, validationErrors, nil
	}

	// 2. Admin şifresini hash'le
	hashedPassword, err := utils.HashPassword(req.AdminPassword)
	if err != nil {
		return nil, nil, fmt.Errorf("şifre hash'lenemedi: %v", err)
	}

	// 3. Hastane ve admin kullanıcıyı oluştur (Repository üzerinden)
	response, err := s.hospitalRepo.CreateHospitalWithAdmin(req, hashedPassword)
	if err != nil {
		return nil, nil, fmt.Errorf("hastane kaydı başarısız: %v", err)
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
