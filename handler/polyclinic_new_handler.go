package handler

import (
	"net/http"
	"strconv"

	"hospital-platform/model"
	"hospital-platform/service"
	"hospital-platform/utils"

	"github.com/labstack/echo/v4"
)

// PolyclinicNewHandler yeni poliklinik HTTP isteklerini yönetir
type PolyclinicNewHandler struct {
	polyclinicService *service.PolyclinicService
}

// NewPolyclinicNewHandler yeni bir poliklinik handler'ı oluşturur
func NewPolyclinicNewHandler() *PolyclinicNewHandler {
	return &PolyclinicNewHandler{
		polyclinicService: service.NewPolyclinicService(),
	}
}

// GetPolyclinicTypes master data poliklinik türlerini getirir
// @Summary Poliklinik türleri listesi
// @Description Master data'dan poliklinik türlerini getirir (dropdown için)
// @Tags Polyclinic
// @Produce json
// @Success 200 {array} model.PolyclinicType
// @Failure 500 {object} map[string]interface{}
// @Router /polyclinic-types [get]
func (h *PolyclinicNewHandler) GetPolyclinicTypes(c echo.Context) error {
	types, err := h.polyclinicService.GetAllPolyclinicTypes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Poliklinik türleri getirilirken hata oluştu",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": types,
	})
}

// AddPolyclinicToHospital hastaneye yeni poliklinik ekler
// @Summary Hastaneye poliklinik ekle
// @Description Seçilen poliklinik türünü hastaneye ekler
// @Tags Polyclinic
// @Accept json
// @Produce json
// @Param body body model.AddPolyclinicRequest true "Poliklinik ekleme verisi"
// @Success 201 {object} model.HospitalPolyclinic
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hospital/polyclinics [post]
func (h *PolyclinicNewHandler) AddPolyclinicToHospital(c echo.Context) error {
	// JWT token'dan hospital ID al
	hospitalID, err := h.getHospitalIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Geçersiz token",
		})
	}

	var req model.AddPolyclinicRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   "Geçersiz istek formatı",
			"details": err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   "Veri doğrulama hatası",
			"details": err.Error(),
		})
	}

	polyclinic, err := h.polyclinicService.AddPolyclinicToHospital(&req, hospitalID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Poliklinik başarıyla eklendi",
		"data":    polyclinic,
	})
}

// GetHospitalPolyclinics hastane polikliniklerini temel bilgilerle getirir
// @Summary Hastane polikliniklerini getir
// @Description Hastaneye ait poliklinikleri personel sayılarıyla listeler
// @Tags Polyclinic
// @Produce json
// @Success 200 {array} model.HospitalPolyclinicSummary
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hospital/polyclinics [get]
func (h *PolyclinicNewHandler) GetHospitalPolyclinics(c echo.Context) error {
	// JWT token'dan hospital ID al
	hospitalID, err := h.getHospitalIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Geçersiz token",
		})
	}

	polyclinics, err := h.polyclinicService.GetHospitalPolyclinics(hospitalID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Poliklinikler getirilirken hata oluştu",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": polyclinics,
	})
}

// UpdateHospitalPolyclinic hastane poliklinik günceller
// @Summary Hastane poliklinik güncelle
// @Description Hastane poliklinik bilgilerini günceller
// @Tags Polyclinic
// @Accept json
// @Produce json
// @Param id path int true "Poliklinik ID"
// @Param body body model.UpdatePolyclinicRequest true "Güncelleme verisi"
// @Success 200 {object} model.HospitalPolyclinic
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hospital/polyclinics/{id} [put]
func (h *PolyclinicNewHandler) UpdateHospitalPolyclinic(c echo.Context) error {
	// JWT token'dan hospital ID al
	hospitalID, err := h.getHospitalIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Geçersiz token",
		})
	}

	// Path parametresi
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Geçersiz poliklinik ID",
		})
	}

	var req model.UpdatePolyclinicRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   "Geçersiz istek formatı",
			"details": err.Error(),
		})
	}

	polyclinic, err := h.polyclinicService.UpdateHospitalPolyclinic(uint(id), &req, hospitalID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Poliklinik başarıyla güncellendi",
		"data":    polyclinic,
	})
}

// DeleteHospitalPolyclinic hastane poliklinik siler
// @Summary Hastane poliklinik sil
// @Description Hastane poliklinik siler
// @Tags Polyclinic
// @Produce json
// @Param id path int true "Poliklinik ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hospital/polyclinics/{id} [delete]
func (h *PolyclinicNewHandler) DeleteHospitalPolyclinic(c echo.Context) error {
	// JWT token'dan hospital ID al
	hospitalID, err := h.getHospitalIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Geçersiz token",
		})
	}

	// Path parametresi
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Geçersiz poliklinik ID",
		})
	}

	err = h.polyclinicService.DeleteHospitalPolyclinic(uint(id), hospitalID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Poliklinik başarıyla silindi",
	})
}

// getHospitalIDFromToken JWT token'dan hospital ID'yi çıkarır
func (h *PolyclinicNewHandler) getHospitalIDFromToken(c echo.Context) (uint, error) {
	// Authorization header'dan token al
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "Authorization header eksik")
	}

	// "Bearer " prefix'ini kaldır
	tokenString := authHeader[7:] // "Bearer " 7 karakter

	// Token'ı doğrula
	_, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "Geçersiz token")
	}

	// TODO: User ID'den Hospital ID'yi çek
	// Şimdilik basit implementation - gerçekte user service'den çekilmeli
	return 1, nil // Geçici olarak 1 döndürüyoruz
}
