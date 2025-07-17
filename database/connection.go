package database

import (
	"fmt"
	"hospital-platform/config"
	"hospital-platform/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	// PostgreSQL connection string - mevcut .env ayarlarını kullan
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.GetEnv("DB_HOST", "localhost"),
		config.GetEnv("DB_USER", "postgres"),
		config.GetEnv("DB_PASSWORD", "Admin..200101"),
		config.GetEnv("DB_NAME", "hospital_db"),
		config.GetEnv("DB_PORT", "5432"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanına bağlanılamadı:", err)
	}

	fmt.Println("Veritabanı bağlantısı başarılı!")

	// Migration yapmadan önce problemli tabloları temizle
	fmt.Println("Mevcut tabloları temizleniyor...")
	dropTables()

	// Run migrations for all models
	err = DB.AutoMigrate(
		// Master data tables (location)
		&model.Province{},
		&model.District{},

		// Master data tables (job system)
		&model.JobGroup{},
		&model.JobTitle{},

		// Master data tables (polyclinic types)
		&model.PolyclinicType{},

		// Main business tables
		&model.Hospital{},
		&model.User{},
		&model.HospitalPolyclinic{},
		&model.Staff{},

		// Legacy tables (backward compatibility)
		&model.Polyclinic{},
		&model.LoginRequest{},
		&model.ResetPasswordRequest{},
		&model.ResetPasswordConfirm{},
	)

	if err != nil {
		log.Fatal("Migration işlemi başarısız:", err)
	}

	fmt.Println("Migration tamamlandı!")

	// Seed master data
	seedMasterData()
}

// dropTables removes problematic tables to allow clean migration
func dropTables() {
	// Önce foreign key constraint'leri olan tabloları sil
	DB.Migrator().DropTable(&model.User{})
	DB.Migrator().DropTable(&model.Staff{})
	DB.Migrator().DropTable(&model.HospitalPolyclinic{})
	DB.Migrator().DropTable(&model.Hospital{})
	DB.Migrator().DropTable(&model.Polyclinic{})

	fmt.Println("Eski tablolar temizlendi.")
}

// seedMasterData populates master data tables
func seedMasterData() {
	fmt.Println("Seeding master data...")

	// 1. Seed Provinces
	provinces := []model.Province{
		{Name: "İstanbul"},
		{Name: "Ankara"},
		{Name: "İzmir"},
		{Name: "Bursa"},
		{Name: "Antalya"},
	}

	for _, province := range provinces {
		var existing model.Province
		result := DB.Where("name = ?", province.Name).First(&existing)
		if result.Error != nil {
			DB.Create(&province)
			fmt.Printf("Created province: %s\n", province.Name)
		}
	}

	// 2. Seed Districts
	districts := []model.District{
		// İstanbul
		{ProvinceID: 1, Name: "Beşiktaş"},
		{ProvinceID: 1, Name: "Kadıköy"},
		{ProvinceID: 1, Name: "Şişli"},
		// Ankara
		{ProvinceID: 2, Name: "Çankaya"},
		{ProvinceID: 2, Name: "Keçiören"},
		// İzmir
		{ProvinceID: 3, Name: "Bornova"},
		{ProvinceID: 3, Name: "Karşıyaka"},
	}

	for _, district := range districts {
		var existing model.District
		result := DB.Where("province_id = ? AND name = ?", district.ProvinceID, district.Name).First(&existing)
		if result.Error != nil {
			DB.Create(&district)
			fmt.Printf("Created district: %s\n", district.Name)
		}
	}

	// 3. Seed Job Groups
	jobGroups := []model.JobGroup{
		{Name: "Doktor"},
		{Name: "İdari Personel"},
		{Name: "Hizmet Personeli"},
	}

	for _, jobGroup := range jobGroups {
		var existing model.JobGroup
		result := DB.Where("name = ?", jobGroup.Name).First(&existing)
		if result.Error != nil {
			DB.Create(&jobGroup)
			fmt.Printf("Created job group: %s\n", jobGroup.Name)
		}
	}

	// 4. Seed Job Titles
	jobTitles := []model.JobTitle{
		// Doktor grubu
		{JobGroupID: 1, Name: "Asistan Doktor", IsUnique: false},
		{JobGroupID: 1, Name: "Uzman Doktor", IsUnique: false},
		{JobGroupID: 1, Name: "Başhekim", IsUnique: true}, // Hastanede sadece 1 tane
		// İdari Personel
		{JobGroupID: 2, Name: "İnsan Kaynakları Uzmanı", IsUnique: false},
		{JobGroupID: 2, Name: "Müdür", IsUnique: false},
		{JobGroupID: 2, Name: "Genel Koordinatör", IsUnique: true},
		// Hizmet Personeli
		{JobGroupID: 3, Name: "Hemşire", IsUnique: false},
		{JobGroupID: 3, Name: "Danışman", IsUnique: false},
		{JobGroupID: 3, Name: "Temizlik Personeli", IsUnique: false},
		{JobGroupID: 3, Name: "Güvenlik", IsUnique: false},
	}

	for _, jobTitle := range jobTitles {
		var existing model.JobTitle
		result := DB.Where("job_group_id = ? AND name = ?", jobTitle.JobGroupID, jobTitle.Name).First(&existing)
		if result.Error != nil {
			DB.Create(&jobTitle)
			fmt.Printf("Created job title: %s\n", jobTitle.Name)
		}
	}

	// 5. Seed Polyclinic Types
	polyclinicTypes := []model.PolyclinicType{
		{Name: "Kardiyoloji", Description: "Kalp ve damar hastalıkları"},
		{Name: "Nöroloji", Description: "Sinir sistemi hastalıkları"},
		{Name: "Göz", Description: "Göz hastalıkları"},
		{Name: "Kulak Burun Boğaz", Description: "KBB hastalıkları"},
		{Name: "Ortopedi", Description: "Kemik ve eklem hastalıkları"},
		{Name: "Dahiliye", Description: "İç hastalıkları"},
		{Name: "Çocuk", Description: "Çocuk hastalıkları"},
		{Name: "Kadın Doğum", Description: "Kadın hastalıkları ve doğum"},
	}

	for _, polyclinicType := range polyclinicTypes {
		var existing model.PolyclinicType
		result := DB.Where("name = ?", polyclinicType.Name).First(&existing)
		if result.Error != nil {
			DB.Create(&polyclinicType)
			fmt.Printf("Created polyclinic type: %s\n", polyclinicType.Name)
		}
	}

	fmt.Println("Master data seeding completed!")
}
