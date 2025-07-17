package handler

import (
	"net/http"
	"strconv"

	"hospital-platform/model"
	"hospital-platform/service"
	"hospital-platform/utils"

	"github.com/labstack/echo/v4"
)

// StaffHandler - Personel yönetimi ile ilgili HTTP isteklerini karşılayan controller katmanı
// API endpoint'lerini handle eder ve servis katmanı ile iletişim kurar
type StaffHandler struct {
	staffService *service.StaffService // İş mantığını yöneten servis
}

// NewStaffHandler - Personel handler'ının yeni bir instance'ını oluşturur
// Servis bağımlılığını enjekte eder ve kullanıma hazır hale getirir
func NewStaffHandler() *StaffHandler {
	return &StaffHandler{
		staffService: service.NewStaffService(),
	}
}

// ==================== PERSONEL YÖNETİM API'LERİ ====================

// CreateStaff - Hastaneye yeni personel ekleme işlemini gerçekleştirir
// @Summary Personel ekle
// @Description Hastaneye yeni personel ekler
// @Tags Staff
// @Accept json
// @Produce json
// @Param body body model.CreateStaffRequest true "Personel ekleme verisi"
// @Success 201 {object} model.Staff
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hospital/staff [post]
func (h *StaffHandler) CreateStaff(c echo.Context) error {
	// JWT token'dan hospital ID al
	hospitalID, err := h.getHospitalIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Geçersiz token",
		})
	}

	var req model.CreateStaffRequest
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

	staff, validationErrors, err := h.staffService.CreateStaff(&req, hospitalID)

	// Validation hataları
	if len(validationErrors) > 0 {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"error":             "Veri doğrulama hataları",
			"validation_errors": validationErrors,
		})
	}

	// Service hataları
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Personel başarıyla eklendi",
		"data":    staff,
	})
}

// GetStaffByID ID'ye göre personel getirir
// @Summary Personel detayları
// @Description ID'ye göre personel detaylarını getirir
// @Tags Staff
// @Produce json
// @Param id path int true "Personel ID"
// @Success 200 {object} model.Staff
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hospital/staff/{id} [get]
func (h *StaffHandler) GetStaffByID(c echo.Context) error {
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
			"error": "Geçersiz personel ID",
		})
	}

	staff, err := h.staffService.GetStaffByID(uint(id), hospitalID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": staff,
	})
}

// UpdateStaff personel bilgilerini günceller
// @Summary Personel güncelle
// @Description Personel bilgilerini günceller
// @Tags Staff
// @Accept json
// @Produce json
// @Param id path int true "Personel ID"
// @Param body body model.UpdateStaffRequest true "Güncelleme verisi"
// @Success 200 {object} model.Staff
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hospital/staff/{id} [put]
func (h *StaffHandler) UpdateStaff(c echo.Context) error {
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
			"error": "Geçersiz personel ID",
		})
	}

	var req model.UpdateStaffRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   "Geçersiz istek formatı",
			"details": err.Error(),
		})
	}

	staff, validationErrors, err := h.staffService.UpdateStaff(uint(id), &req, hospitalID)

	// Validation hataları
	if len(validationErrors) > 0 {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"error":             "Veri doğrulama hataları",
			"validation_errors": validationErrors,
		})
	}

	// Service hataları
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Personel başarıyla güncellendi",
		"data":    staff,
	})
}

// DeleteStaff personeli siler
// @Summary Personel sil
// @Description Personeli siler
// @Tags Staff
// @Produce json
// @Param id path int true "Personel ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hospital/staff/{id} [delete]
func (h *StaffHandler) DeleteStaff(c echo.Context) error {
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
			"error": "Geçersiz personel ID",
		})
	}

	err = h.staffService.DeleteStaff(uint(id), hospitalID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Personel başarıyla silindi",
	})
}

// ==================== LİSTELEME VE FİLTRELEME ====================

// GetStaffList sayfalandırılmış personel listesi getirir
// @Summary Personel listesi
// @Description Hastane personellerini sayfalandırılmış ve filtreli olarak getirir
// @Tags Staff
// @Accept json
// @Produce json
// @Param body body model.StaffListRequest true "Listeleme ve filtreleme verisi"
// @Success 200 {object} model.StaffListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hospital/staff/list [post]
func (h *StaffHandler) GetStaffList(c echo.Context) error {
	// JWT token'dan hospital ID al
	hospitalID, err := h.getHospitalIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Geçersiz token",
		})
	}

	var req model.StaffListRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   "Geçersiz istek formatı",
			"details": err.Error(),
		})
	}

	response, err := h.staffService.GetStaffList(&req, hospitalID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

// ==================== MASTER DATA ====================

// GetJobGroups meslek gruplarını getirir
// @Summary Meslek grupları
// @Description Tüm meslek gruplarını getirir
// @Tags Staff
// @Produce json
// @Success 200 {array} model.JobGroup
// @Failure 500 {object} map[string]interface{}
// @Router /job-groups [get]
func (h *StaffHandler) GetJobGroups(c echo.Context) error {
	jobGroups, err := h.staffService.GetJobGroups()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Meslek grupları getirilirken hata oluştu",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": jobGroups,
	})
}

// GetJobTitlesByGroup meslek grubuna göre unvanları getirir
// @Summary Meslek grubuna göre unvanlar
// @Description Seçilen meslek grubuna ait unvanları getirir
// @Tags Staff
// @Produce json
// @Param job_group_id path int true "Meslek grubu ID"
// @Success 200 {array} model.JobTitle
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /job-groups/{job_group_id}/titles [get]
func (h *StaffHandler) GetJobTitlesByGroup(c echo.Context) error {
	// Path parametresi
	jobGroupIDParam := c.Param("job_group_id")
	jobGroupID, err := strconv.ParseUint(jobGroupIDParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Geçersiz meslek grubu ID",
		})
	}

	jobTitles, err := h.staffService.GetJobTitlesByGroup(uint(jobGroupID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Unvanlar getirilirken hata oluştu",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": jobTitles,
	})
}

// ==================== HELPER METHODS ====================

// getHospitalIDFromToken JWT token'dan hospital ID'yi çıkarır
func (h *StaffHandler) getHospitalIDFromToken(c echo.Context) (uint, error) {
	// Authorization header'dan token al
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "Authorization header eksik")
	}

	// "Bearer " prefix'ini kaldır
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "Geçersiz token formatı")
	}
	tokenString := authHeader[7:]

	// Token'ı doğrula
	_, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "Geçersiz token")
	}

	// TODO: User ID'den Hospital ID'yi çek
	// Şimdilik basit implementation - gerçekte user service'den çekilmeli
	return 1, nil // Geçici olarak 1 döndürüyoruz
}
