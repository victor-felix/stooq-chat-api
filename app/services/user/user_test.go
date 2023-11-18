package user_test

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/victor-felix/chat-service/app/models"
	"github.com/victor-felix/chat-service/app/services/user"
)

func TestCreate_User(t *testing.T) {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	t.Run("should return error when email is empty", func(t *testing.T) {
		service := user.NewUserService(nil, "", 0, log)
		payload := &models.User{
			Email: "",
			Password: "password",
			UserName: "username",
		}
		user, err := service.Create(context.Background(), payload)
		
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.Equal(t, "<REQUIRED_FIELD_MISSING> email is required", err.Error())
	})

	t.Run("should return error when password is empty", func(t *testing.T) {
		service := user.NewUserService(nil, "", 0, log)
		payload := &models.User{
			Email: "email",
			Password: "",
			UserName: "username",
		}
		user, err := service.Create(context.Background(), payload)
		
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.Equal(t, "<REQUIRED_FIELD_MISSING> password is required", err.Error())
	})

	t.Run("should return error when username is empty", func(t *testing.T) {
		service := user.NewUserService(nil, "", 0, log)
		payload := &models.User{
			Email: "email",
			Password: "password",
			UserName: "",
		}
		user, err := service.Create(context.Background(), payload)
		
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.Equal(t, "<REQUIRED_FIELD_MISSING> username is required", err.Error())
	})

	t.Run("should return error when email is already in use", func(t *testing.T) {
		userStorageMock := &user.StorageMock{
			FindOneByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
				return &models.User{
					Email: "email",
					Password: "password",
					UserName: "username",
				}, nil
			},
		}

		service := user.NewUserService(userStorageMock, "", 0, log)
		payload := &models.User{
			Email: "email",
			Password: "password",
			UserName: "username",
		}

		user, err := service.Create(context.Background(), payload)

		assert.Nil(t, user)
		assert.Error(t, err)
		assert.Equal(t, "<EMAIL_DUPLICATED> email already exists", err.Error())
	});

	t.Run("should return error when find one by email returns error", func(t *testing.T) {
		userStorageMock := &user.StorageMock{
			FindOneByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
				return nil, assert.AnError
			},
		}

		service := user.NewUserService(userStorageMock, "", 0, log)
		payload := &models.User{
			Email: "email",
			Password: "password",
			UserName: "username",
		}

		user, err := service.Create(context.Background(), payload)

		assert.Nil(t, user)
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should return auth response payload when user storage returns user", func(t *testing.T) {
		userStorageMock := &user.StorageMock{
			FindOneByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
				return nil, nil
			},
			InsertFn: func(ctx context.Context, user *models.User) (*models.User, error) {
				return &models.User{
					ID: "id",
					Email: "email",
					Password: "password",
					UserName: "username",
				}, nil
			},
		}

		service := user.NewUserService(userStorageMock, "", 0, log)
		payload := &models.User{
			Email: "email",
			Password: "password",
			UserName: "username",
		}

		user, err := service.Create(context.Background(), payload)

		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.NotNil(t, user.Token)
	})
}