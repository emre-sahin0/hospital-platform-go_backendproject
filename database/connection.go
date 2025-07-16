package database

import (
	"fmt"
	"log"
	"os"

	"hospital-platform/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("Veritabanına bağlanılamadı:", err)
	}

	// 🔽 Otomatik tablo oluşturma
	err = DB.AutoMigrate(
		&model.User{},
		&model.Polyclinic{}, // 🔥 Burayı ekledik
	)

	if err != nil {
		log.Panic("Migration işlemi başarısız:", err)
	}

	fmt.Println("Veritabanı bağlantısı başarılı")
}
