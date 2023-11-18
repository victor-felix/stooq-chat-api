package message_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victor-felix/chat-service/app/models"
	"github.com/victor-felix/chat-service/app/services/message"
)

func TestCreate_Message(t *testing.T) {
	t.Run("should return error when room id is empty", func(t *testing.T) {
		service := message.NewMessageService(nil)
		message := &models.Message{
			RoomID: "",
			UserID: "user-id",
			Content: "content",
		}
		message, err := service.Create(context.Background(), message)

		assert.Nil(t, message)
		assert.Error(t, err)
		assert.Equal(t, "<REQUIRED_FIELD_MISSING> room_id is required", err.Error())
	})

	t.Run("should return error when user id is empty", func(t *testing.T) {
		service := message.NewMessageService(nil)
		message := &models.Message{
			RoomID: "room-id",
			UserID: "",
			Content: "content",
		}
		message, err := service.Create(context.Background(), message)

		assert.Nil(t, message)
		assert.Error(t, err)
		assert.Equal(t, "<REQUIRED_FIELD_MISSING> user_id is required", err.Error())
	})

	t.Run("should return error when content is empty", func(t *testing.T) {
		service := message.NewMessageService(nil)
		message := &models.Message{
			RoomID: "room-id",
			UserID: "user-id",
			Content: "",
		}
		message, err := service.Create(context.Background(), message)

		assert.Nil(t, message)
		assert.Error(t, err)
		assert.Equal(t, "<REQUIRED_FIELD_MISSING> content is required", err.Error())
	})

	t.Run("should return error when message store return error", func(t *testing.T) {
		messageStorageMock := &message.StorageMock{
			InsertFn: func(ctx context.Context, message *models.Message) (*models.Message, error) {
				return nil, assert.AnError
			},
		}

		service := message.NewMessageService(messageStorageMock)
		message := &models.Message{
			RoomID: "room-id",
			UserID: "user-id",
			Content: "content",
		}
		message, err := service.Create(context.Background(), message)

		assert.Nil(t, message)
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should return message when message store returns message", func(t *testing.T) {
		messageStorageMock := &message.StorageMock{
			InsertFn: func(ctx context.Context, message *models.Message) (*models.Message, error) {
				return message, nil
			},
		}

		service := message.NewMessageService(messageStorageMock)
		message := &models.Message{
			RoomID: "room-id",
			UserID: "user-id",
			Content: "content",
		}
		message, err := service.Create(context.Background(), message)

		assert.NotNil(t, message)
		assert.NoError(t, err)
	});
}
