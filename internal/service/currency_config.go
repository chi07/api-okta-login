package service

import (
	"context"

	"github.com/chi07/api-okta-login/internal/http/request"
)

type CurrencyConfigService struct {
	currencyConfigRepo CurrencyConfigRepo
}

func NewCurrencyConfigService(currencyConfigRepo CurrencyConfigRepo) *CurrencyConfigService {
	return &CurrencyConfigService{
		currencyConfigRepo: currencyConfigRepo,
	}
}

func (s *CurrencyConfigService) BulkUpdateExclusiveCurrency(ctx context.Context, currencies []*request.ExclusiveCurrency) error {
	return s.currencyConfigRepo.BulkUpdateExclusiveCurrency(ctx, currencies)
}
