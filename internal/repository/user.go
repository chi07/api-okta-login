package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/chi07/api-okta-login/internal/model"
)

type User struct {
	db *mongo.Database
}

func NewUser(db *mongo.Database) *User {
	return &User{db: db}
}

func (repo *User) Collection() *mongo.Collection {
	return repo.db.Collection("users")
}

func (repo *User) Get(ctx context.Context, id string) (*model.User, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("user id = %s is not valid", id))
	}
	filter := bson.D{{Key: "_id", Value: userID}}
	user := &model.User{}

	err = repo.Collection().FindOne(ctx, filter).Decode(&user)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.Wrap(err, "cannot get users from database")
	}

	return user, nil
}

func (repo *User) GetByOktaID(ctx context.Context, oktaID string) (*model.User, error) {
	filter := bson.D{{Key: "oktaID", Value: oktaID}}
	var user *model.User

	err := repo.Collection().FindOne(ctx, filter).Decode(&user)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.Wrap(err, "cannot get user from database")
	}

	return user, nil
}

func (repo *User) Create(ctx context.Context, u *model.User) (*model.User, error) {
	document, err := repo.Collection().InsertOne(ctx, u)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get insert user into database")
	}

	filter := bson.D{{Key: "_id", Value: document.InsertedID}}
	createdRecord := repo.Collection().FindOne(ctx, filter)

	var user *model.User
	err = createdRecord.Decode(&user)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get user from database")
	}

	return user, nil
}
