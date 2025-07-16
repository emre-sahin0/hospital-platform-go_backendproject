package handler

import (
	"net/http"
	"strconv"

	"hospital-platform/model"
	"hospital-platform/service"

	"github.com/labstack/echo/v4"
)

// @Summary Yeni poliklinik oluştur
// @Description Yeni bir poliklinik ekler.
// @Tags Polyclinic
// @Accept json
// @Produce json
// @Param body body model.Polyclinic true "Poliklinik verisi"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /polyclinics [post]
func CreatePolyclinic(c echo.Context) error {
	var poly model.Polyclinic

	if err := c.Bind(&poly); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Geçersiz veri"})
	}

	if err := service.CreatePolyclinicService(&poly); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Poliklinik başarıyla eklendi",
		"data":    poly,
	})
}

// @Summary Tüm poliklinikleri getir
// @Description Sistemdeki tüm poliklinikleri listeler.
// @Tags Polyclinic
// @Produce json
// @Success 200 {array} model.Polyclinic
// @Failure 500 {object} map[string]interface{}
// @Router /polyclinics [get]
func GetAllPolyclinics(c echo.Context) error {
	data, err := service.GetAllPolyclinicsService()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Listeleme hatası"})
	}
	return c.JSON(http.StatusOK, data)
}

// @Summary Poliklinik sil
// @Description ID'ye göre poliklinik kaydını siler.
// @Tags Polyclinic
// @Produce json
// @Param id path int true "Poliklinik ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /polyclinics/{id} [delete]
func DeletePolyclinic(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := service.DeletePolyclinicService(uint(id)); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Silinemedi"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Poliklinik başarıyla silindi"})
}

// @Summary Poliklinik bilgilerini güncelle
// @Tags Polyclinic
// @Accept json
// @Produce json
// @Param id path int true "Poliklinik ID"
// @Param body body model.Polyclinic true "Yeni poliklinik verisi"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /polyclinics/{id} [put]
func UpdatePolyclinic(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var poly model.Polyclinic
	if err := c.Bind(&poly); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Geçersiz veri"})
	}
	if err := service.UpdatePolyclinicService(uint(id), &poly); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Poliklinik güncellendi", "data": poly})
}
