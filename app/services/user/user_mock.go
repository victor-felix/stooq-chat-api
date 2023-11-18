package user

import (
	"context"

	"github.com/victor-felix/chat-service/app/models"
)

type ServiceMock struct {
	CreateFn func(ctx context.Context, user *models.User) (*models.AuthResponsePayload, error)
	CreateFnCallCount int
	FindByEmailFn func(ctx context.Context, email string) (*models.User, error)
	FindByEmailFnCallCount int
	UserValidateFn func(ctx context.Context, user *models.User) error
	UserValidateFnCallCount int
}

func (sm *ServiceMock) Create(ctx context.Context, user *models.User) (*models.AuthResponsePayload, error) {
	sm.CreateFnCallCount++
	return sm.CreateFn(ctx, user)
}

func (sm *ServiceMock) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	sm.FindByEmailFnCallCount++
	return sm.FindByEmailFn(ctx, email)
}

func (sm *ServiceMock) UserValidate(ctx context.Context, user *models.User) error {
	sm.UserValidateFnCallCount++
	return sm.UserValidateFn(ctx, user)
}

type StorageMock struct {
	InsertFn func(ctx context.Context, user *models.User) (*models.User, error)
	InsertFnCallCount int
	FindOneByEmailFn func(ctx context.Context, email string) (*models.User, error)
	FindOneByEmailFnCallCount int
}

func (sm *StorageMock) Insert(ctx context.Context, user *models.User) (*models.User, error) {
	sm.InsertFnCallCount++
	return sm.InsertFn(ctx, user)
}

func (sm *StorageMock) FindOneByEmail(ctx context.Context, email string) (*models.User, error) {
	sm.FindOneByEmailFnCallCount++
	return sm.FindOneByEmailFn(ctx, email)
}
