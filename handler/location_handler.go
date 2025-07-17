package handler

import (
	"net/http"
	"strconv"

	"hospital-platform/service"

	"github.com/labstack/echo/v4"
)

// LocationHandler lokasyon HTTP isteklerini yönetir
type LocationHandler struct {
	locationService *service.LocationService
}

// NewLocationHandler yeni bir lokasyon handler'ı oluşturur
func NewLocationHandler() *LocationHandler {
	return &LocationHandler{
		locationService: service.NewLocationService(),
	}
}

// GetAllProvinces tüm illeri getirir
// @Summary Tüm illeri getir
// @Description Dropdown için tüm illeri listeler
// @Tags Location
// @Produce json
// @Success 200 {array} model.Province
// @Failure 500 {object} map[string]interface{}
// @Router /provinces [get]
func (h *LocationHandler) GetAllProvinces(c echo.Context) error {
	provinces, err := h.locationService.GetAllProvinces()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "İller getirilirken hata oluştu",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": provinces,
	})
}

// GetDistrictsByProvinceID belirtilen il için ilçeleri getirir
// @Summary İle göre ilçeleri getir
// @Description Seçilen ile ait ilçeleri listeler
// @Tags Location
// @Produce json
// @Param province_id path int true "İl ID"
// @Success 200 {array} model.District
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /provinces/{province_id}/districts [get]
func (h *LocationHandler) GetDistrictsByProvinceID(c echo.Context) error {
	provinceIDParam := c.Param("province_id")
	provinceID, err := strconv.ParseUint(provinceIDParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Geçersiz il ID",
		})
	}

	districts, err := h.locationService.GetDistrictsByProvinceID(uint(provinceID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "İlçeler getirilirken hata oluştu",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": districts,
	})
}
