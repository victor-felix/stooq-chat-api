package message

import (
	"context"

	"github.com/victor-felix/chat-service/app/models"
)

type ServiceMock struct {
	FindByRoomIdFn func(ctx context.Context, roomID string) ([]*models.Message, error)
	FindByRoomIdFnCallCount int
	CreateFn func(ctx context.Context, message *models.Message) (*models.Message, error)
	CreateFnCallCount int
}

func (sm *ServiceMock) FindByRoomId(ctx context.Context, roomID string) ([]*models.Message, error) {
	sm.FindByRoomIdFnCallCount++
	return sm.FindByRoomIdFn(ctx, roomID)
}

func (sm *ServiceMock) Create(ctx context.Context, message *models.Message) (*models.Message, error) {
	sm.CreateFnCallCount++
	return sm.CreateFn(ctx, message)
}

type StorageMock struct {
	FindAllByRoomIdFn func(ctx context.Context, roomID string) ([]*models.Message, error)
	FindAllByRoomIdFnCallCount int
	InsertFn func(ctx context.Context, message *models.Message) (*models.Message, error)
	InsertFnCallCount int
}

func (sm *StorageMock) FindAllByRoomId(ctx context.Context, roomID string) ([]*models.Message, error) {
	sm.FindAllByRoomIdFnCallCount++
	return sm.FindAllByRoomIdFn(ctx, roomID)
}

func (sm *StorageMock) Insert(ctx context.Context, message *models.Message) (*models.Message, error) {
	sm.InsertFnCallCount++
	return sm.InsertFn(ctx, message)
}
