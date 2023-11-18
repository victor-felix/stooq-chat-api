package message

import (
	"context"

	"github.com/victor-felix/chat-service/app/models"
)

type MessageService struct {
	messageStorage models.MessageStorage
}

func NewMessageService(messageStorage models.MessageStorage) models.MessageService {
	return &MessageService{
		messageStorage: messageStorage,
	}
}

func (ms *MessageService) FindByRoomId(ctx context.Context, roomID string) ([]*models.Message, error) {
	return ms.messageStorage.FindAllByRoomId(ctx, roomID)
}

func (ms *MessageService) Create(ctx context.Context, message *models.Message) (*models.Message, error) {
	if err := message.Validate(); err != nil {
		return nil, err
	}

	if _, err := ms.messageStorage.Insert(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

func (ms *MessageService) SaveMessageByWebSocket(roomID, userID, content string) bool {
	ctx := context.Background()

	message := &models.Message{
		RoomID:  roomID,
		UserID:  userID,
		Content: content,
	}

	if err := message.Validate(); err != nil {
		return false
	}

	if _, err := ms.messageStorage.Insert(ctx, message); err != nil {
		return false
	}

	return true
}

func (ms *MessageService) SaveMessageByAmqp(message models.Message) bool {
	ctx := context.Background()

	if _, err := ms.messageStorage.Insert(ctx, &message); err != nil {
		return false
	}

	return true
}
