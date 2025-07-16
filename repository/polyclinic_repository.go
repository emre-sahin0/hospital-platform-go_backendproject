package repository

import (
	"hospital-platform/database"
	"hospital-platform/model"
)

// Yeni poliklinik oluştur
func CreatePolyclinic(polyclinic *model.Polyclinic) error {
	return database.DB.Create(polyclinic).Error
}

// Tüm poliklinikleri getir
func GetAllPolyclinics() ([]model.Polyclinic, error) {
	var polyclinics []model.Polyclinic
	err := database.DB.Find(&polyclinics).Error
	return polyclinics, err
}

// ID ile poliklinik getir

// Poliklinik güncelle

// Poliklinik sil
func DeletePolyclinic(id uint) error {
	return database.DB.Delete(&model.Polyclinic{}, id).Error
}

func GetPolyclinicByID(id uint) (*model.Polyclinic, error) {
	var poly model.Polyclinic
	if err := database.DB.First(&poly, id).Error; err != nil {
		return nil, err
	}
	return &poly, nil
}

func UpdatePolyclinic(poly *model.Polyclinic) error {
	return database.DB.Save(poly).Error
}
