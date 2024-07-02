package service

import (
	"context"

	"github.com/chi07/api-okta-login/internal/model"
)

type LoginHistoryRepo interface {
	Create(ctx context.Context, c *model.LoginHistory) (*model.LoginHistory, error)
}

type UserRepo interface {
	Create(ctx context.Context, u *model.User) (*model.User, error)
	Get(ctx context.Context, id string) (*model.User, error)
	GetByOktaID(ctx context.Context, oktaID string) (*model.User, error)
}
