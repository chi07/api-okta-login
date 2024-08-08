package service

import (
	"context"

	"github.com/chi07/api-okta-login/internal/http/request"
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

type CurrencyConfigRepo interface {
	BulkUpdateExclusiveCurrency(ctx context.Context, currencies []*request.ExclusiveCurrency) error
}
