package room_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victor-felix/chat-service/app/models"
	"github.com/victor-felix/chat-service/app/services/room"
)

func TestCreate_Room(t *testing.T){
	t.Run("should return error when name is empty", func(t *testing.T) {
		service := room.NewRoomService(nil)
		room := &models.Room{
			Name: "",
		}
		room, err := service.Create(context.Background(), room)

		assert.Nil(t, room)
		assert.Error(t, err)
		assert.Equal(t, "<REQUIRED_FIELD_MISSING> name is required", err.Error())
	})

	t.Run("should return error when room store return error", func(t *testing.T) {
		roomStorageMock := &room.StorageMock{
			InsertFn: func(ctx context.Context, room *models.Room) (*models.Room, error) {
				return nil, assert.AnError
			},
		}
		service := room.NewRoomService(roomStorageMock)
		room := &models.Room{
			Name: "room-name",
		}
		room, err := service.Create(context.Background(), room)

		assert.Nil(t, room)
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should return room when room store return room", func(t *testing.T) {
		roomStorageMock := &room.StorageMock{
			InsertFn: func(ctx context.Context, room *models.Room) (*models.Room, error) {
				return &models.Room{
					ID: "",
					Name: "room-name",
				}, nil
			},
		}
		service := room.NewRoomService(roomStorageMock)
		room := &models.Room{
			Name: "room-name",
		}
		room, err := service.Create(context.Background(), room)

		assert.Nil(t, err)
		assert.Equal(t, "", room.ID)
		assert.Equal(t, "room-name", room.Name)
	})
}