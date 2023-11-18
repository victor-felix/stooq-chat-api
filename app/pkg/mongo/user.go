package mongo

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/victor-felix/chat-service/app/errors"
	"github.com/victor-felix/chat-service/app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStorage struct {
	collection *mongo.Collection
	log zerolog.Logger
}

func NewUserStorage(db *mongo.Database, log zerolog.Logger) models.UserStorage {
	return &UserStorage{
		collection: db.Collection("users"),
		log: log,
	}
}

func (us *UserStorage) FindOneByEmail(ctx context.Context, email string) (*models.User, error) {
	result := &models.User{}
	if err := us.collection.FindOne(ctx, map[string]string{"email": email}).Decode(result); err != nil {
		if IsNotFoundError(err) {
			return nil, errors.NewErrNotFound(models.ErrorUserNotFound).WithMessage("user not found")
		}
		us.log.Error().Msg(err.Error())
		return nil, err
	}

	return result, nil
}

func (us *UserStorage) Insert(ctx context.Context, user *models.User) (*models.User, error) {
	user.ID = primitive.NewObjectID().Hex()
	if _, err := us.collection.InsertOne(ctx, user); err != nil {
		if IsDuplicatedError(err) {
			return nil, errors.NewErrInvalidArgument(models.ErrorEmailDuplicated).WithMessage("email already exists")
		}
		us.log.Error().Msg(err.Error())
		return nil, err
	}

	return user, nil
}
