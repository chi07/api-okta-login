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

	//// save to sessions
	//sess, err := session.Get("session", c)
	//if err != nil {
	//	return err
	//}
	//sess.Options = &sessions.Options{
	//	Path:     "/",
	//	MaxAge:   86400 * 7,
	//	HttpOnly: true,
	//}
	//sess.Values["accessToken"] = result["accessToken"]
	//if err := sess.Save(c.Request(), c.Response()); err != nil {
	//	return err
	//}

	return c.JSON(http.StatusOK, result)
}

func (h *LoginHandler) Logout(c echo.Context) error {
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
