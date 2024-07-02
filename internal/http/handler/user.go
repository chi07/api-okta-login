package handler

import (
	"github.com/chi07/api-okta-login/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"net/http"
)

type LoginHandler struct {
	config *config.Config
	log    zerolog.Logger
}

func NewLoginHandler(config *config.Config, log zerolog.Logger) *LoginHandler {
	return &LoginHandler{
		config: config,
		log:    log,
	}
}

func (h *LoginHandler) LoginWithOkta(c echo.Context) error {
	idToken := c.FormValue("id_token")

	token, err := verifyToken(idToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid_token", "error_description": err.Error()})
	}

	// Token is valid, proceed with your application logic
	// For example, create a session, return user info, etc.

	return c.JSON(http.StatusOK, token)
}

func verifyToken(idToken string) (map[string]interface{}, error) {
	return nil, nil
}
