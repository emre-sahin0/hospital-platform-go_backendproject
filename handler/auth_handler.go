package handler

import (
	"fmt"
	"net/http"

	"hospital-platform/model"
	"hospital-platform/service"

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
// @Description Email ve şifre ile kullanıcı girişi yapılır
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

	token, err := service.Login(credentials.Email, credentials.Password)
	if err != nil {
		fmt.Println("Login hatası:", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}
