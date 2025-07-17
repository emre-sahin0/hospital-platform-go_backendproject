package handler

import (
	"fmt"
	"math/rand"
	"net/http"

	"hospital-platform/database"
	"hospital-platform/model"
	"hospital-platform/repository"
	"hospital-platform/utils"

	"github.com/labstack/echo/v4"
)

// ResetPasswordRequest godoc
// @Summary Şifre sıfırlama kodu gönder
// @Description Telefon numarasına geçici şifre sıfırlama kodu üretir
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.ResetPasswordRequest true "Telefon numarası"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /reset-password/request [post]
// Şifre sıfırlama isteği: Telefon numarasına geçici kod üret
func ResetPasswordRequestHandler(c echo.Context) error {
	var request struct {
		Phone string `json:"phone"`
	}
	if err := c.Bind(&request); err != nil {
		fmt.Println("Bind hatası:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Geçersiz istek"})
	}

	// Kullanıcı veritabanında var mı?
	user, err := repository.NewUserRepository().GetByPhone(request.Phone)
	if err != nil {
		fmt.Println("Kullanıcı sorgusu hatası:", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Sunucu hatası"})
	}
	if user == nil {
		fmt.Println("Telefon numarası sistemde kayıtlı değil:", request.Phone)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Telefon numarası bulunamadı"})
	}

	// 6 haneli kod üret
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Redis'e kaydet (key: phone, value: code)
	err = database.SetResetCode(request.Phone, code)
	if err != nil {
		fmt.Println("Redis'e kod yazılamadı:", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Kod üretilemedi"})
	}

	fmt.Printf("Şifre sıfırlama kodu gönderildi: %s -> %s\n", request.Phone, code)

	// Kullanıcıya kod üretildi mesajı
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Kod üretildi",
		"code":    code, // Gerçekte SMS gönderilir ama test için döndük
	})
}

// ResetPasswordConfirm godoc
// @Summary Kodla şifreyi sıfırla
// @Description Kod doğruysa yeni şifre veritabanına yazılır
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.ResetPasswordConfirm true "Yeni şifre bilgileri"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /reset-password/confirm [post]
func ResetPasswordConfirm(c echo.Context) error {
	var request model.ResetPasswordConfirm

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Geçersiz veri"})
	}

	// Redis’ten kodu al
	storedCode, err := database.RedisClient.Get(database.Ctx, request.Phone).Result()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Kod süresi dolmuş veya bulunamadı"})
	}

	fmt.Println("✔️ Beklenen kod:", storedCode)
	fmt.Println("📲 Girilen kod:", request.Code)

	if request.Code != storedCode {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Kod eşleşmiyor"})
	}

	if request.NewPassword != request.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Şifreler uyuşmuyor"})
	}

	// Kullanıcıyı bul
	var user model.User
	if err := database.DB.Where("phone = ?", request.Phone).First(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Kullanıcı bulunamadı"})
	}

	// Yeni şifreyi hashle
	hashed, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Hashleme hatası"})
	}

	user.Password = hashed
	if err := database.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Şifre güncellenemedi"})
	}

	// Kod Redis’ten silinir
	database.RedisClient.Del(database.Ctx, request.Phone)

	fmt.Println("🔐 Şifre başarıyla güncellendi")
	return c.JSON(http.StatusOK, echo.Map{"message": "Şifre başarıyla güncellendi"})
}
