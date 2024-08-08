package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/chi07/api-okta-login/internal/config"
	"github.com/chi07/api-okta-login/internal/http/request"
)

type CurrencyConfigHandler struct {
	config          *config.Config
	log             zerolog.Logger
	validator       *validator.Validate
	currencyService CurrencyConfigService
}

func NewCurrencyConfigHandler(config *config.Config, log zerolog.Logger, validator *validator.Validate, currencyService CurrencyConfigService) *CurrencyConfigHandler {
	return &CurrencyConfigHandler{
		config:          config,
		log:             log,
		validator:       validator,
		currencyService: currencyService,
	}
}

func (h *CurrencyConfigHandler) BulkUpdateExclusiveCurrency(c echo.Context) error {
	var req request.BulkExclusiveCurrency

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

	err := h.currencyService.BulkUpdateExclusiveCurrency(c.Request().Context(), req.Currencies)
	if err != nil {
		h.log.Error().Err(err).Msg("failed to update exclusive currency")
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
			"code":    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}
