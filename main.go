// @title Hastane Takip API
// @version 1.0
// @description Bu API VatanSoft staj projesi için yazılmıştır.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token için "Bearer " prefix'i ile birlikte token'ı girin. Örnek: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

package main

import (
	"hospital-platform/config"
	"hospital-platform/database"
	"hospital-platform/handler"
	"hospital-platform/utils" // Middleware'ler için

	_ "hospital-platform/docs" // Swagger docs

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// CustomValidator struct for Echo validation
type CustomValidator struct {
	validator *validator.Validate
}

// Validate method for Echo validation interface
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// Ortam değişkenlerini yükle
	config.LoadEnv() // .env'den verileri çeksin

	// Veritabanına bağlan
	database.ConnectDB()
	database.ConnectRedis()

	// Echo başlat
	e := echo.New()

	// Validator middleware'i kur
	e.Validator = &CustomValidator{validator: validator.New()}

	// API dokümantasyonu için Swagger UI endpoint'i
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Tüm handler'ları initialize et - dependency injection
	hospitalHandler := handler.NewHospitalHandler()           // Hastane yönetimi
	locationHandler := handler.NewLocationHandler()           // İl/İlçe dropdown'ları
	polyclinicNewHandler := handler.NewPolyclinicNewHandler() // Poliklinik yönetimi
	staffHandler := handler.NewStaffHandler()                 // Personel yönetimi

	// ========== 🌍 AÇIK ERİŞİM ROTALARİ (Middleware Yok) ==========

	// Kimlik doğrulama - herkes erişebilir
	e.POST("/login", handler.Login)
	e.POST("/register", handler.Register)
	e.POST("/reset-password/request", handler.ResetPasswordRequestHandler)
	e.POST("/reset-password/confirm", handler.ResetPasswordConfirm)

	// Hastane kaydı - herkes erişebilir
	e.POST("/hospital/register", hospitalHandler.RegisterHospital)

	// Master data - herkes erişebilir (dropdown'lar için)
	e.GET("/provinces", locationHandler.GetAllProvinces)
	e.GET("/provinces/:province_id/districts", locationHandler.GetDistrictsByProvinceID)
	e.GET("/polyclinic-types", polyclinicNewHandler.GetPolyclinicTypes)
	e.GET("/job-groups", staffHandler.GetJobGroups)
	e.GET("/job-groups/:job_group_id/titles", staffHandler.GetJobTitlesByGroup)

	// ========== 🔐 KORUNMUŞ ERİŞİM ROTALARİ (JWT Gerekli) ==========

	// JWT middleware'i olan grup oluştur
	protected := e.Group("")
	protected.Use(utils.JWTAuthMiddleware())

	// Hastane bilgileri - login olan herkes görebilir
	protected.GET("/hospital/:id", hospitalHandler.GetHospitalByID)

	// ========== 👀 OKUMA İZNİ GEREKLİ (Hem Yetkili Hem Çalışan) ==========

	// Okuma izni olan grup oluştur
	readAccess := protected.Group("")
	readAccess.Use(utils.RequirePermission(utils.READ))

	// Poliklinik görüntüleme - hem yetkili hem çalışan
	readAccess.GET("/hospital/polyclinics", polyclinicNewHandler.GetHospitalPolyclinics)

	// Personel görüntüleme - hem yetkili hem çalışan
	readAccess.GET("/hospital/staff/:id", staffHandler.GetStaffByID)
	readAccess.POST("/hospital/staff/list", staffHandler.GetStaffList) // Filtreleme dahil

	// ========== 🔒 YÖNETİCİ İZNİ GEREKLİ (Sadece Yetkili) ==========

	// Admin izni olan grup oluştur
	adminAccess := protected.Group("")
	adminAccess.Use(utils.RequirePermission(utils.ADMIN))

	// Poliklinik yönetimi - sadece yetkili
	adminAccess.POST("/hospital/polyclinics", polyclinicNewHandler.AddPolyclinicToHospital)
	adminAccess.PUT("/hospital/polyclinics/:id", polyclinicNewHandler.UpdateHospitalPolyclinic)
	adminAccess.DELETE("/hospital/polyclinics/:id", polyclinicNewHandler.DeleteHospitalPolyclinic)

	// Personel yönetimi - sadece yetkili
	adminAccess.POST("/hospital/staff", staffHandler.CreateStaff)
	adminAccess.PUT("/hospital/staff/:id", staffHandler.UpdateStaff)
	adminAccess.DELETE("/hospital/staff/:id", staffHandler.DeleteStaff)

	// Alt kullanıcı yönetimi - sadece yetkili
	adminAccess.POST("/hospital/users", handler.CreateSubUser)
	adminAccess.GET("/hospital/users", handler.GetSubUsers)
	adminAccess.PUT("/hospital/users/:id", handler.UpdateSubUser)
	adminAccess.DELETE("/hospital/users/:id", handler.DeleteSubUser)

	// ========== POLYCLINIC ROUTES (Legacy - Geriye Uyumluluk) ==========
	// Legacy polyclinic endpoints
	e.GET("/polyclinics", handler.GetAllPolyclinics)
	e.POST("/polyclinics", handler.CreatePolyclinic)
	e.PUT("/polyclinics/:id", handler.UpdatePolyclinic)
	e.DELETE("/polyclinics/:id", handler.DeletePolyclinic)

	// Sunucuyu başlat
	e.Logger.Fatal(e.Start(":8080"))
}
