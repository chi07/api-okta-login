package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/chi07/api-okta-login/internal/config"
	"github.com/chi07/api-okta-login/internal/http/request"
)

type LoginHandler struct {
	config      *config.Config
	log         zerolog.Logger
	validator   *validator.Validate
	oktaService OktaService
}

func NewLoginHandler(config *config.Config, log zerolog.Logger, validator *validator.Validate, oktaService OktaService) *LoginHandler {
	return &LoginHandler{
		config:      config,
		log:         log,
		validator:   validator,
		oktaService: oktaService,
	}
}

func (h *LoginHandler) GetMyIP(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"ip": c.RealIP(),
	})
}

func (h *LoginHandler) Logout(c echo.Context) error {
	// @TODO: xoá access token trong session như hôm qua đi
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "logout",
	})
}

func (h *LoginHandler) LoginWithOkta(c echo.Context) error {
	var req request.LoginRequest

	if err := c.Bind(&req); err != nil {
		h.log.Error().Err(err).Msg("failed to bind request")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"code":    http.StatusBadRequest,
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.Error().Err(err).Msg("failed to validate request")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"code":    http.StatusBadRequest,
		})
	}

	result, err := h.oktaService.Login(c.Request().Context(), req.OktaToken)
	if err != nil {
		h.log.Error().Err(err).Msg("failed to login")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"code":    http.StatusBadRequest,
		})
	}

	return c.JSON(http.StatusOK, result)
}
