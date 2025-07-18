package utils

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// JWTClaims - Token'dan çıkardığımız kullanıcı bilgileri
type JWTClaims struct {
	UserID     uint   `json:"user_id"`
	HospitalID uint   `json:"hospital_id"`
	Role       string `json:"role"`
	Username   string `json:"username"`
}

// PermissionLevel - Yetki seviyeleri enum'u
type PermissionLevel int

const (
	READ  PermissionLevel = iota // Sadece okuma (çalışan)
	WRITE                        // Okuma + yazma (çalışan+)
	ADMIN                        // Her şey (yetkili)
)

// JWTAuthMiddleware - JWT token'ı doğrular ve context'e kullanıcı bilgilerini ekler
// Her korumalı endpoint'te bu middleware çalışır
func JWTAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorization header'ını kontrol et
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error":   "Yetkilendirme hatası",
					"message": "Authorization header eksik",
				})
			}

			// "Bearer " prefix'ini kontrol et
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error":   "Yetkilendirme hatası",
					"message": "Geçersiz token formatı",
				})
			}

			// Token'ı çıkar
			tokenString := authHeader[7:] // "Bearer " kısmını kes

			// Token'ı doğrula ve parse et
			claims, err := ValidateJWT(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error":   "Yetkilendirme hatası",
					"message": "Geçersiz veya süresi dolmuş token",
				})
			}

			// Claims'i context'e ekle - diğer handler'lar kullanabilsin
			c.Set("user_id", claims["user_id"])
			c.Set("hospital_id", claims["hospital_id"])
			c.Set("role", claims["role"])
			c.Set("username", claims["username"])

			return next(c)
		}
	}
}

// RequirePermission - Belirli yetki seviyesi gerektiren endpoint'ler için middleware
// Örnek: RequirePermission(ADMIN) - sadece yetkili kullanıcılar erişebilir
func RequirePermission(requiredLevel PermissionLevel) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// JWT middleware'den role bilgisini al
			roleInterface := c.Get("role")
			if roleInterface == nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error":   "Yetkilendirme hatası",
					"message": "Rol bilgisi bulunamadı",
				})
			}

			role, ok := roleInterface.(string)
			if !ok {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"error":   "Sistem hatası",
					"message": "Rol bilgisi işlenemedi",
				})
			}

			// Kullanıcının yetki seviyesini belirle
			userLevel := getUserPermissionLevel(role)

			// Gerekli yetki seviyesi ile karşılaştır
			if userLevel < requiredLevel {
				return c.JSON(http.StatusForbidden, echo.Map{
					"error":   "Yetkisiz erişim",
					"message": "Bu işlem için yeterli yetkiniz yok",
					"details": map[string]string{
						"required": getPermissionName(requiredLevel),
						"current":  getPermissionName(userLevel),
					},
				})
			}

			return next(c)
		}
	}
}

// RequireRole - Belirli rol gerektiren endpoint'ler için middleware
// Örnek: RequireRole("yetkili") - sadece yetkili kullanıcılar
func RequireRole(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// JWT middleware'den role bilgisini al
			roleInterface := c.Get("role")
			if roleInterface == nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error":   "Yetkilendirme hatası",
					"message": "Rol bilgisi bulunamadı",
				})
			}

			role, ok := roleInterface.(string)
			if !ok {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"error":   "Sistem hatası",
					"message": "Rol bilgisi işlenemedi",
				})
			}

			// Rol kontrolü
			if role != requiredRole {
				return c.JSON(http.StatusForbidden, echo.Map{
					"error":   "Yetkisiz erişim",
					"message": "Bu işlem için '" + requiredRole + "' rolü gerekli",
					"details": map[string]string{
						"required": requiredRole,
						"current":  role,
					},
				})
			}

			return next(c)
		}
	}
}

// getUserPermissionLevel - Kullanıcı rolünden yetki seviyesini belirler
func getUserPermissionLevel(role string) PermissionLevel {
	switch role {
	case "yetkili":
		return ADMIN // Yetkili kullanıcı - her şeyi yapabilir
	case "çalışan":
		return READ // Çalışan - sadece okuma
	default:
		return READ // Bilinmeyen rol - güvenlik için sadece okuma
	}
}

// getPermissionName - Yetki seviyesinin ismini döndürür (hata mesajları için)
func getPermissionName(level PermissionLevel) string {
	switch level {
	case READ:
		return "okuma"
	case WRITE:
		return "yazma"
	case ADMIN:
		return "yönetici"
	default:
		return "bilinmeyen"
	}
}

// =============== HELPER FUNCTIONS ===============

// GetUserIDFromContext - Context'ten user ID'yi çıkarır
func GetUserIDFromContext(c echo.Context) (uint, bool) {
	userIDInterface := c.Get("user_id")
	if userIDInterface == nil {
		return 0, false
	}

	// interface{} to float64 (JWT'den gelen sayılar float64 olur)
	userIDFloat, ok := userIDInterface.(float64)
	if !ok {
		return 0, false
	}

	return uint(userIDFloat), true
}

// GetHospitalIDFromContext - Context'ten hospital ID'yi çıkarır
func GetHospitalIDFromContext(c echo.Context) (uint, bool) {
	hospitalIDInterface := c.Get("hospital_id")
	if hospitalIDInterface == nil {
		return 0, false
	}

	// interface{} to float64 (JWT'den gelen sayılar float64 olur)
	hospitalIDFloat, ok := hospitalIDInterface.(float64)
	if !ok {
		return 0, false
	}

	return uint(hospitalIDFloat), true
}

// GetRoleFromContext - Context'ten rol bilgisini çıkarır
func GetRoleFromContext(c echo.Context) (string, bool) {
	roleInterface := c.Get("role")
	if roleInterface == nil {
		return "", false
	}

	role, ok := roleInterface.(string)
	return role, ok
}

// GetUsernameFromContext - Context'ten kullanıcı adını çıkarır
func GetUsernameFromContext(c echo.Context) (string, bool) {
	usernameInterface := c.Get("username")
	if usernameInterface == nil {
		return "", false
	}

	username, ok := usernameInterface.(string)
	return username, ok
}
