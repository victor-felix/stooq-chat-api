package room

import (
	"context"

	"github.com/victor-felix/chat-service/app/models"
)

type ServiceMock struct {
	FindAllFn func(ctx context.Context) ([]*models.Room, error)
	FindAllFnCallCount int
	FindByIDFn func(ctx context.Context, id string) (*models.Room, error)
	FindByIDFnCallCount int
	CreateFn func(ctx context.Context, room *models.Room) (*models.Room, error)
	CreateFnCallCount int
}

func (sm *ServiceMock) FindAll(ctx context.Context) ([]*models.Room, error) {
	sm.FindAllFnCallCount++
	return sm.FindAllFn(ctx)
}

func (sm *ServiceMock) FindByID(ctx context.Context, id string) (*models.Room, error) {
	sm.FindByIDFnCallCount++
	return sm.FindByIDFn(ctx, id)
}

func (sm *ServiceMock) Create(ctx context.Context, room *models.Room) (*models.Room, error) {
	sm.CreateFnCallCount++
	return sm.CreateFn(ctx, room)
}

type StorageMock struct {
	FindAllFn func(ctx context.Context) ([]*models.Room, error)
	FindAllFnCallCount int
	FindOneByIDFn func(ctx context.Context, id string) (*models.Room, error)
	FindOneByIDFnCallCount int
	InsertFn func(ctx context.Context, room *models.Room) (*models.Room, error)
	InsertFnCallCount int
}

func (sm *StorageMock) FindAll(ctx context.Context) ([]*models.Room, error) {
	sm.FindAllFnCallCount++
	return sm.FindAllFn(ctx)
}

func (sm *StorageMock) FindOneByID(ctx context.Context, id string) (*models.Room, error) {
	sm.FindOneByIDFnCallCount++
	return sm.FindOneByIDFn(ctx, id)
}

func (sm *StorageMock) Insert(ctx context.Context, room *models.Room) (*models.Room, error) {
	sm.InsertFnCallCount++
	return sm.InsertFn(ctx, room)
}
