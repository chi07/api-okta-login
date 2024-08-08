package repository

import (
	"context"
	"fmt"

	"github.com/chi07/api-okta-login/internal/http/request"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CurrencyConfig struct {
	db *mongo.Database
}

func NewCurrencyConfig(db *mongo.Database) *CurrencyConfig {
	return &CurrencyConfig{db: db}
}

func (repo *CurrencyConfig) Collection() *mongo.Collection {
	return repo.db.Collection("currency_configs")
}

func (repo *CurrencyConfig) BulkUpdateExclusiveCurrency(ctx context.Context, currencies []*request.ExclusiveCurrency) error {
	for _, p := range currencies {
		filter := bson.M{"currency": p.Currency}
		update := bson.M{"$set": bson.M{"is_exclusive": p.IsExclusive}}

		_, err := repo.db.Collection("currency_configs").UpdateOne(ctx, filter, update)
		if err != nil {
			return fmt.Errorf("failed to update document: %v", err)
		}
	}

	return nil
}
