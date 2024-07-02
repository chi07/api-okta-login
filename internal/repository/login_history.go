package repository

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/chi07/api-okta-login/internal/model"
)

type LoginHistory struct {
	db *mongo.Database
}

func NewLoginHistory(db *mongo.Database) *LoginHistory {
	return &LoginHistory{db: db}
}

func (repo *LoginHistory) Collection() *mongo.Collection {
	return repo.db.Collection("loginHistories")
}

func (repo *LoginHistory) Create(ctx context.Context, lh *model.LoginHistory) (*model.LoginHistory, error) {
	collection := repo.db.Collection("loginHistory")
	document, err := collection.InsertOne(ctx, lh)
	if err != nil {
		return nil, errors.Wrap(err, "cannot insert loginHistory into database")
	}

	filter := bson.D{{Key: "_id", Value: document.InsertedID}}
	createdRecord := collection.FindOne(ctx, filter)

	var loginHistory *model.LoginHistory
	err = createdRecord.Decode(&loginHistory)
	if err != nil {
		return nil, errors.Wrap(err, "cannot Decode loginHistory from database")
	}

	return loginHistory, nil
}
