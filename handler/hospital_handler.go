package handler

import (
	"net/http"
	"strconv"

	"hospital-platform/model"
	"hospital-platform/service"

	"github.com/labstack/echo/v4"
)

// HospitalHandler hastane HTTP isteklerini yönetir
type HospitalHandler struct {
	hospitalService *service.HospitalService
}

// NewHospitalHandler yeni bir hastane handler'ı oluşturur
func NewHospitalHandler() *HospitalHandler {
	return &HospitalHandler{
		hospitalService: service.NewHospitalService(),
	}
}

// RegisterHospital yeni hastane ve admin kullanıcı kaydı yapar
// @Summary Hastane kaydı
// @Description Yeni hastane ve admin kullanıcı kaydı yapar
// @Tags Hospital
// @Accept json
// @Produce json
// @Param body body model.HospitalRegistrationRequest true "Hastane kayıt verisi"
// @Success 201 {object} model.HospitalRegistrationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Router /hospital/register [post]
func (h *HospitalHandler) RegisterHospital(c echo.Context) error {
	var req model.HospitalRegistrationRequest

	// İstek body'sini parse et
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   "Geçersiz istek formatı",
			"details": err.Error(),
		})
	}

	// İstek alanlarını doğrula
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   "Veri doğrulama hatası",
			"details": err.Error(),
		})
	}

	// Hastane kaydını yap
	response, validationErrors, err := h.hospitalService.RegisterHospital(&req)

	// Doğrulama hatalarını işle
	if len(validationErrors) > 0 {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"error":             "Veri doğrulama hataları",
			"validation_errors": validationErrors,
		})
	}

	// Servis hatalarını işle
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   "Hastane kaydı sırasında hata oluştu",
			"details": err.Error(),
		})
	}

	// Başarılı yanıt
	return c.JSON(http.StatusCreated, response)
}

// GetHospitalByID ID'ye göre hastane bilgilerini getirir
// @Summary Hastane bilgilerini getir
// @Description ID'ye göre hastane bilgilerini getirir
// @Tags Hospital
// @Produce json
// @Param id path int true "Hastane ID"
// @Success 200 {object} model.Hospital
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /hospital/{id} [get]
func (h *HospitalHandler) GetHospitalByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Geçersiz hastane ID",
		})
	}

	hospital, err := h.hospitalService.GetHospitalByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Hastane bulunamadı",
		})
	}

	return c.JSON(http.StatusOK, hospital)
}
