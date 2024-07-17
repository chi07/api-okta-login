package middleware

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strings"

	"github.com/chi07/api-okta-login/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type AuthMiddleWare struct {
	store     *sessions.CookieStore
	jwtConfig *config.JWTConfig
}

func NewAuthMiddleWare(store *sessions.CookieStore, jwtConfig *config.JWTConfig) *AuthMiddleWare {
	return &AuthMiddleWare{store: store, jwtConfig: jwtConfig}
}

func (m *AuthMiddleWare) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("auth middleware")
		jwtToken, claims, err := m.secureEndpoint(c)
		if !jwtToken.Valid || err != nil {
			return errors.New("invalid token")
		}

		// Get user info from the claims
		fmt.Println(claims)

		// @TODO Kiểm tra accessToken có ở trong sessions ko, nếu không có thì trả về lỗi token không hợp lệ

		// Execute the handler
		err = next(c)
		if err != nil {
			c.Error(err)
		}

		return nil
	}
}

// getBearerToken extracts the Bearer token from the Authorization header
func (m *AuthMiddleWare) getBearerToken(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization header format")
	}

	return parts[1], nil
}

// decodeJWT decodes and validates the JWT token
func (m *AuthMiddleWare) decodeJWT(tokenString string) (*jwt.Token, *jwt.MapClaims, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify the token signature with your secret key or public key
		return []byte("oYUMa8rT8EHUfOzZ0U0Ul5FzMBgM0DO4"), nil
	})
	return token, claims, err
}

// secureEndpoint is a sample endpoint that requires a Bearer token
func (m *AuthMiddleWare) secureEndpoint(c echo.Context) (*jwt.Token, *jwt.MapClaims, error) {
	tokenString, err := m.getBearerToken(c)
	if err != nil {
		return nil, nil, err
	}

	return m.decodeJWT(tokenString)
}
