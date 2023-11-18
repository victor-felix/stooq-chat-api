package models

import (
	"context"
	"strings"
	"time"

	"github.com/victor-felix/chat-service/app/errors"
)

type Message struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`
	RoomID string `json:"room_id,omitempty" bson:"room_id,omitempty"`
	Room *Room `json:"room,omitempty" bson:"room,omitempty"`
	Content string `json:"content,omitempty" bson:"content,omitempty"`
	UserID string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	User *User `json:"user,omitempty" bson:"user,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

func (m *Message) Validate() error {
	var field string

	field = strings.TrimSpace(m.Content)
	if len(field) == 0 {
		return errors.NewErrInvalidArgument(ErrorRequiredFieldMissing).WithMessage("content is required")
	}
	
	field = strings.TrimSpace(m.RoomID)
	if len(field) == 0 {
		return errors.NewErrInvalidArgument(ErrorRequiredFieldMissing).WithMessage("room_id is required")
	}

	field = strings.TrimSpace(m.UserID)
	if len(field) == 0 {
		return errors.NewErrInvalidArgument(ErrorRequiredFieldMissing).WithMessage("user_id is required")
	}

	return nil
}

type MessageStorage interface {
	Insert(ctx context.Context, message *Message) (*Message, error)
	FindAllByRoomId(ctx context.Context, roomID string) ([]*Message, error)
}

type MessageService interface {
	SaveMessageByWebSocket(roomID, userID, content string) bool
	SaveMessageByAmqp(message Message) bool
	Create(ctx context.Context, message *Message) (*Message, error)
	FindByRoomId(ctx context.Context, roomID string) ([]*Message, error)
}