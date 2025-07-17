// @title Hastane Takip API
// @version 1.0
// @description Bu API VatanSoft staj giriş projesi için yazılmıştır.
// @host localhost:8080
// @BasePath /

package main

import (
	"hospital-platform/config"
	"hospital-platform/database"
	"hospital-platform/handler"

	_ "hospital-platform/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {

	config.LoadEnv()

	// Veritabanına bağlan
	database.ConnectDB()
	database.ConnectRedis()

	// Echo başlat
	e := echo.New()

	// API dokümantasyonu için Swagger UI endpoint'i
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Tüm handler'ları initialize et - dependency injection
	hospitalHandler := handler.NewHospitalHandler()           // Hastane yönetimi
	locationHandler := handler.NewLocationHandler()           // İl/İlçe dropdown'ları
	polyclinicNewHandler := handler.NewPolyclinicNewHandler() // Poliklinik yönetimi
	staffHandler := handler.NewStaffHandler()                 // Personel yönetimi

	e.POST("/login", handler.Login)
	e.POST("/register", handler.Register)

	// Şifre sıfırlama süreçleri
	e.POST("/reset-password/request", handler.ResetPasswordRequestHandler)
	e.POST("/reset-password/confirm", handler.ResetPasswordConfirm)

	// ========== HASTANE YÖNETİM ROTALARİ ==========
	// Yeni hastane kayıt sistemi ve hastane bilgileri
	e.POST("/hospital/register", hospitalHandler.RegisterHospital) // Hastane kaydı + admin kullanıcı oluşturma
	e.GET("/hospital/:id", hospitalHandler.GetHospitalByID)        // Hastane detaylarını görüntüleme

	// ========== COĞRAFİ VERİ ROTALARİ ==========
	// İl ve ilçe dropdown listeleri için master data
	e.GET("/provinces", locationHandler.GetAllProvinces)                                 // 81 ili listele
	e.GET("/provinces/:province_id/districts", locationHandler.GetDistrictsByProvinceID) // Seçili ile ait ilçeleri listele

	// ========== POLİKLİNİK YÖNETİM ROTALARİ ==========
	// Master data - sistemde tanımlı poliklinik türleri
	e.GET("/polyclinic-types", polyclinicNewHandler.GetPolyclinicTypes)

	// Hastane bazlı poliklinik yönetimi (JWT token gerekli)
	e.POST("/hospital/polyclinics", polyclinicNewHandler.AddPolyclinicToHospital)        // Hastaneye poliklinik ekle
	e.GET("/hospital/polyclinics", polyclinicNewHandler.GetHospitalPolyclinics)          // Hastanedeki poliklinikleri listele
	e.PUT("/hospital/polyclinics/:id", polyclinicNewHandler.UpdateHospitalPolyclinic)    // Poliklinik güncelle
	e.DELETE("/hospital/polyclinics/:id", polyclinicNewHandler.DeleteHospitalPolyclinic) // Poliklinik sil

	// ========== PERSONEL YÖNETİM ROTALARİ ==========
	// Master data - meslek grupları ve unvanlar
	e.GET("/job-groups", staffHandler.GetJobGroups)                             // Tüm meslek gruplarını listele (Doktor, Hemşire, vb.)
	e.GET("/job-groups/:job_group_id/titles", staffHandler.GetJobTitlesByGroup) // Seçilen gruba ait unvanları listele

	// Personel CRUD işlemleri (JWT token gerekli)
	e.POST("/hospital/staff", staffHandler.CreateStaff)       // Yeni personel ekle
	e.GET("/hospital/staff/:id", staffHandler.GetStaffByID)   // Personel detaylarını görüntüle
	e.PUT("/hospital/staff/:id", staffHandler.UpdateStaff)    // Personel bilgilerini güncelle
	e.DELETE("/hospital/staff/:id", staffHandler.DeleteStaff) // Personeli sil

	// Personel listeleme ve filtreleme (JWT token gerekli)
	e.POST("/hospital/staff/list", staffHandler.GetStaffList) // Sayfalandırılmış ve filtreli personel listesi

	// ========== POLYCLINIC ROUTES (Legacy - Geriye Uyumluluk) ==========
	// Legacy polyclinic endpoints
	e.GET("/polyclinics", handler.GetAllPolyclinics)
	e.POST("/polyclinics", handler.CreatePolyclinic)
	e.PUT("/polyclinics/:id", handler.UpdatePolyclinic)
	e.DELETE("/polyclinics/:id", handler.DeletePolyclinic)

	// Sunucuyu başlat
	e.Logger.Fatal(e.Start(":8080"))
}
