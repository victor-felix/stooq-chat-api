package models

import (
	"context"
	"strings"
	"time"

	"github.com/victor-felix/chat-service/app/errors"
)

type Room struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

func (r *Room) Validate() error {
	field := strings.TrimSpace(r.Name)
	if len(field) == 0 {
		return errors.NewErrInvalidArgument(ErrorRequiredFieldMissing).WithMessage("name is required")
	}

	return nil
}

type RoomStorage interface {
	Insert(ctx context.Context, room *Room) (*Room, error)
	FindAll(ctx context.Context) ([]*Room, error)
	FindOneByID(ctx context.Context, id string) (*Room, error)
}

type RoomService interface {
	Create(ctx context.Context, room *Room) (*Room, error)
	FindAll(ctx context.Context) ([]*Room, error)
	FindByID(ctx context.Context, id string) (*Room, error)
}