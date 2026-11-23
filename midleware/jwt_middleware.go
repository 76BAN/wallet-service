package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"wallets_service/database"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

// Claims berisi informasi yang disimpan dalam JWT.
type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTMiddleware menyimpan public key untuk memvalidasi token.
type JWTMiddleware struct {
	PublicKey *rsa.PublicKey
	RedisDB   *database.RedisService
}

// NewJWTMiddleware membaca public key dari file PEM.
func NewJWTMiddleware(
	publicKeyPath string,
	redisDB *database.RedisService) (*JWTMiddleware, error) {
	data, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca public key: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("format public key PEM tidak valid")
	}

	// Mendukung format public key PKIX dan PKCS#1.
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		publicKey, ok := parsedKey.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("public key bukan RSA")
		}
		return &JWTMiddleware{
			PublicKey: publicKey,
			RedisDB:   redisDB}, nil
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("gagal memproses public key RSA: %w", err)
	}

	return &JWTMiddleware{
		PublicKey: publicKey,
		RedisDB:   redisDB,
	}, nil
}

// Authenticate memvalidasi JWT sebelum handler dijalankan.
func (m *JWTMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Token tidak ditemukan",
			})
		}

		const prefix = "Bearer "
		if len(authHeader) <= len(prefix) ||
			authHeader[:len(prefix)] != prefix {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Format token tidak valid",
			})
		}

		tokenString := authHeader[len(prefix):]

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (any, error) {
				if token.Method != jwt.SigningMethodRS256 {
					return nil, errors.New("algoritma JWT tidak valid")
				}
				return m.PublicKey, nil
			},
		)

		if err != nil || token == nil || !token.Valid {
			fmt.Printf("JWT validation error: %v\n", err)

			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Token tidak valid atau sudah kedaluwarsa",
			})
		}

		if claims.UserID <= 0 {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "UserID tidak ditemukan dalam token",
			})
		}

		// Simpan UserID yang sudah divalidasi ke Echo Context.
		c.Set("userID", claims.UserID)

		exists, err := m.RedisDB.Client.Exists(
			c.Request().Context(),
			"active_token:"+tokenString,
		).Result()

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Gagal mengecek token",
			})
		}
		if exists == 0 {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Token sudah tidak aktif",
			})
		}

		return next(c)
	}
}
