package auth_test

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/victor-felix/chat-service/app/models"
	"github.com/victor-felix/chat-service/app/services/auth"
	"github.com/victor-felix/chat-service/app/services/user"
)


func TestLogin_Auth(t *testing.T) {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	t.Run("should return error when email is empty", func(t *testing.T) {
		service := auth.NewAuthService(nil, "", 0, log)
		payload := &models.AuthPayload{
			Email: "",
			Password: "password",
		}
		user, err := service.Login(context.Background(), payload)
		
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.Equal(t, "<REQUIRED_FIELD_MISSING> email is required", err.Error())
	})

	t.Run("should return error when password is empty", func(t *testing.T) {
		service := auth.NewAuthService(nil, "", 0, log)
		payload := &models.AuthPayload{
			Email: "email",
			Password: "",
		}
		user, err := service.Login(context.Background(), payload)
		
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.Equal(t, "<REQUIRED_FIELD_MISSING> password is required", err.Error())
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		userServiceMock := &user.ServiceMock{
			FindByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
				return nil, nil
			},
		}

		service := auth.NewAuthService(userServiceMock, "", 0, log)
		payload := &models.AuthPayload{
			Email: "email",
			Password: "password",
		}

		user, err := service.Login(context.Background(), payload)

		assert.Nil(t, user)
		assert.Error(t, err)
		assert.Equal(t, "<USER_NOT_FOUND> user not found", err.Error())
	})

	t.Run("should return error wher find by email returns error", func(t *testing.T) {
		userServiceMock := &user.ServiceMock{
			FindByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
				return nil, assert.AnError
			},
		}

		service := auth.NewAuthService(userServiceMock, "", 0, log)
		payload := &models.AuthPayload{
			Email: "email",
			Password: "password",
		}

		user, err := service.Login(context.Background(), payload)

		assert.Nil(t, user)
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		userServiceMock := &user.ServiceMock{
			FindByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
				return &models.User{
					ID: "id",
					UserName: "username",
					Email: "email",
					Password: "password1",
				}, nil
			},
		}

		service := auth.NewAuthService(userServiceMock, "", 0, log)
		payload := &models.AuthPayload{
			Email: "email",
			Password: "password",
		}

		user, err := service.Login(context.Background(), payload)

		assert.Nil(t, user)
		assert.Error(t, err)
		assert.Equal(t, "<INVALID_PASSWORD> invalid password", err.Error())
	})
	
	t.Run("should return token when login is successful", func(t *testing.T) {
		userServiceMock := &user.ServiceMock{
			FindByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
				return &models.User{
					ID: "id",
					UserName: "username",
					Email: "email",
					Password: "$2y$10$8inXIUZLKAwWHSGRrUaMxOkD3h29Qtr7.HBoTl1PvY3P.EyXWs7pG",
				}, nil
			},
		}

		service := auth.NewAuthService(userServiceMock, "secret", 1, log)
		payload := &models.AuthPayload{
			Email: "email",
			Password: "password",
		}

		user, err := service.Login(context.Background(), payload)

		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.NotNil(t, user.Token)
	})
}