package handler

import (
	"context"

	jwtverifier "github.com/okta/okta-jwt-verifier-golang/v2"
)

type OktaService interface {
	Login(ctx context.Context, token string) (map[string]interface{}, error)
	VerifyToken(ctx context.Context, requestToken string) (*jwtverifier.Jwt, bool, error)
}
