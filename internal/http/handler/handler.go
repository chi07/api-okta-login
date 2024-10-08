package handler

import (
	"context"

	"github.com/labstack/echo/v4"

	"github.com/chi07/api-okta-login/internal/http/request"

	jwtverifier "github.com/okta/okta-jwt-verifier-golang/v2"
)

type OktaService interface {
	Login(ctx echo.Context, token string) (map[string]interface{}, error)
	VerifyToken(ctx context.Context, requestToken string) (*jwtverifier.Jwt, bool, error)
	Logout(ctx echo.Context) error
}

type CurrencyConfigService interface {
	BulkUpdateExclusiveCurrency(ctx context.Context, currencies []*request.ExclusiveCurrency) error
}
