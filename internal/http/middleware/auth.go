package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/chi07/api-okta-login/internal/config"
)

type AuthMiddleWare struct {
	store     *sessions.CookieStore
	jwtConfig *config.JWTConfig
	log       zerolog.Logger
}

func NewAuthMiddleWare(store *sessions.CookieStore, jwtConfig *config.JWTConfig, log zerolog.Logger) *AuthMiddleWare {
	return &AuthMiddleWare{store: store, jwtConfig: jwtConfig, log: log}
}

func (m *AuthMiddleWare) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// @TODO Kiểm tra accessToken có hợp lệ không, nếu không hợp lệ thì trả về lỗi token không hợp lệ
		jwtToken, claims, err := m.secureEndpoint(c)
		if !jwtToken.Valid || err != nil {
			m.log.Error().Err(err).Msg("Invalid token")
			return err
		}

		// Get user info from the claims
		var userID, email string
		if sub, ok := (*claims)["sub"].(string); ok {
			userID = sub
		}
		if em, ok := (*claims)["email"].(string); ok {
			email = em
		}

		// @TODO Kiểm tra accessToken có ở trong sessions ko, nếu không có thì trả về lỗi token không hợp lệ
		bearerToken, err := m.getBearerToken(c)
		if err != nil {
			m.log.Error().Err(err).Msg("m.getBearerToken() got err: " + err.Error())
		}

		currentUser, err := m.store.Get(c.Request(), "currentUser")
		if err != nil {
			m.log.Error().Err(err).Msg("m.store.Get() got err: " + err.Error())
		}

		storedAccessToken := currentUser.Values["accessToken"]
		if bearerToken != storedAccessToken {
			m.log.Info().Msg("Invalid token: token not found in session")
			c.Error(err)
			return errors.New("invalid token: token not found in session")
		}

		m.log.Info().Msg("User is authenticated")

		c.Set("userID", userID)
		c.Set("email", email)
		c.Set("accessToken", storedAccessToken)

		// Execute the handler
		if err = next(c); err != nil {
			m.log.Error().Err(err).Msg("next(c) got err: " + err.Error())
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
		return []byte(m.jwtConfig.SecretKey), nil
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
