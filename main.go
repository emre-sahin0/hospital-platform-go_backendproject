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

	_ "hospital-platform/docs" // ✅ Başında alt çizgi (_) olmalı

	"github.com/labstack/echo/v4" // docs klasörü burada import edilir
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// Ortam değişkenlerini yükle

	config.LoadEnv() // .env'den verileri çeksin

	// Veritabanına bağlan
	database.ConnectDB()
	database.ConnectRedis()

	// Echo başlat
	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/login", handler.Login)
	e.POST("/reset-password/request", handler.ResetPasswordRequestHandler)
	e.POST("/reset-password/confirm", handler.ResetPasswordConfirm)

	e.POST("/polyclinics", handler.CreatePolyclinic)
	e.GET("/polyclinics", handler.GetAllPolyclinics)
	e.DELETE("/polyclinics/:id", handler.DeletePolyclinic)
	e.PUT("/polyclinics/:id", handler.UpdatePolyclinic)

	// Kayıt endpoint'i
	e.POST("/register", handler.Register)

	// Sunucuyu başlat
	e.Logger.Fatal(e.Start(":8080"))

}
