package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"hospital-platform/model"
	"hospital-platform/service"
	"hospital-platform/utils"

	"github.com/labstack/echo/v4"
)

// Register godoc
// @Summary Yeni kullanıcı kaydı
// @Description Yeni kullanıcıyı kaydeder
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body model.User true "Kullanıcı bilgileri"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /register [post]
func Register(c echo.Context) error {
	var user model.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Geçersiz veri"})
	}

	if err := service.RegisterUser(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Kayıt başarısız",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Kayıt başarılı"})
}

// Login godoc
// @Summary Kullanıcı girişi
// @Description Email veya telefon numarası ve şifre ile kullanıcı girişi yapılır
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param credentials body model.LoginRequest true "Giriş bilgileri"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /login [post]
func Login(c echo.Context) error {
	var credentials model.LoginRequest
	if err := c.Bind(&credentials); err != nil {
		fmt.Println("Bind hatası:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Geçersiz veri"})
	}
	fmt.Printf("Gelen credentials: %+v\n", credentials)

	token, err := service.Login(credentials.EmailOrPhone, credentials.Password)
	if err != nil {
		fmt.Println("Login hatası:", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

// ==================== ALT KULLANICI YÖNETİMİ ====================

// CreateSubUser godoc
// @Summary Alt kullanıcı ekle
// @Description Yetkili kullanıcı tarafından alt kullanıcı eklenir
// @Tags User Management
// @Accept json
// @Produce json
// @Param user body model.CreateSubUserRequest true "Alt kullanıcı bilgileri"
// @Success 201 {object} model.User
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hospital/users [post]
func CreateSubUser(c echo.Context) error {
	// JWT token'dan kullanıcı bilgilerini al
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Geçersiz token",
		})
	}

	role, ok := utils.GetRoleFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Rol bilgisi bulunamadı",
		})
	}

	// Sadece yetkili kullanıcılar alt kullanıcı ekleyebilir
	if role != "yetkili" {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Bu işlem için yetkili rolü gerekli",
		})
	}

	var req model.CreateSubUserRequest
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

	user, validationErrors, err := service.CreateSubUser(&req, userID)

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
		"message": "Alt kullanıcı başarıyla eklendi",
		"data":    user,
	})
}

// GetSubUsers godoc
// @Summary Alt kullanıcıları listele
// @Description Hastaneye ait alt kullanıcıları listeler
// @Tags User Management
// @Produce json
// @Success 200 {array} model.User
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hospital/users [get]
func GetSubUsers(c echo.Context) error {
	hospitalID, ok := utils.GetHospitalIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Geçersiz token",
		})
	}

	users, err := service.GetSubUsers(hospitalID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": users,
	})
}

// UpdateSubUser godoc
// @Summary Alt kullanıcı güncelle
// @Description Yetkili kullanıcı tarafından alt kullanıcı güncellenir
// @Tags User Management
// @Accept json
// @Produce json
// @Param id path int true "Kullanıcı ID"
// @Param user body model.UpdateSubUserRequest true "Güncelleme verisi"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hospital/users/{id} [put]
func UpdateSubUser(c echo.Context) error {
	// JWT token'dan kullanıcı bilgilerini al
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Geçersiz token",
		})
	}

	role, ok := utils.GetRoleFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Rol bilgisi bulunamadı",
		})
	}

	// Sadece yetkili kullanıcılar alt kullanıcı güncelleyebilir
	if role != "yetkili" {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Bu işlem için yetkili rolü gerekli",
		})
	}

	// Path parametresi
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Geçersiz kullanıcı ID",
		})
	}

	var req model.UpdateSubUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   "Geçersiz istek formatı",
			"details": err.Error(),
		})
	}

	user, validationErrors, err := service.UpdateSubUser(uint(id), &req, userID)

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
		"message": "Alt kullanıcı başarıyla güncellendi",
		"data":    user,
	})
}

// DeleteSubUser godoc
// @Summary Alt kullanıcı sil
// @Description Yetkili kullanıcı tarafından alt kullanıcı silinir
// @Tags User Management
// @Produce json
// @Param id path int true "Kullanıcı ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hospital/users/{id} [delete]
func DeleteSubUser(c echo.Context) error {
	// JWT token'dan kullanıcı bilgilerini al
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Geçersiz token",
		})
	}

	role, ok := utils.GetRoleFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Rol bilgisi bulunamadı",
		})
	}

	// Sadece yetkili kullanıcılar alt kullanıcı silebilir
	if role != "yetkili" {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Bu işlem için yetkili rolü gerekli",
		})
	}

	// Path parametresi
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Geçersiz kullanıcı ID",
		})
	}

	err = service.DeleteSubUser(uint(id), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Alt kullanıcı başarıyla silindi",
	})
}
