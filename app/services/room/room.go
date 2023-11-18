package room

import (
	"context"

	"github.com/victor-felix/chat-service/app/models"
)

type RoomService struct {
	roomStorage models.RoomStorage
}

func NewRoomService(roomStorage models.RoomStorage) models.RoomService {
	return &RoomService{
		roomStorage: roomStorage,
	}
}

func (rs *RoomService) Create(ctx context.Context, room *models.Room) (*models.Room, error) {
	if err := room.Validate(); err != nil {
		return nil, err
	}

	if _, err := rs.roomStorage.Insert(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}

func (rs *RoomService) FindAll(ctx context.Context) ([]*models.Room, error) {
	return rs.roomStorage.FindAll(ctx)
}

func (rs *RoomService) FindByID(ctx context.Context, id string) (*models.Room, error) {
	return rs.roomStorage.FindOneByID(ctx, id)
}
