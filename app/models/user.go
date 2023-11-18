package models

import (
	"context"
	"strings"

	"github.com/victor-felix/chat-service/app/errors"
)

type User struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`
	UserName string `json:"user_name,omitempty" bson:"user_name,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}

func (u *User) Validate() error {
	var field string

	field = strings.TrimSpace(u.Email)
	if len(field) == 0 {
		return errors.NewErrInvalidArgument(ErrorRequiredFieldMissing).WithMessage("email is required")
	}

	field = strings.TrimSpace(u.Password)
	if len(field) == 0 {
		return errors.NewErrInvalidArgument(ErrorRequiredFieldMissing).WithMessage("password is required")
	}

	field = strings.TrimSpace(u.UserName)
	if len(field) == 0 {
		return errors.NewErrInvalidArgument(ErrorRequiredFieldMissing).WithMessage("username is required")
	}

	return nil
}

type UserStorage interface {
	Insert(ctx context.Context, user *User) (*User, error)
	FindOneByEmail(ctx context.Context, email string) (*User, error)
}

type UserService interface {
	Create(ctx context.Context, user *User) (*AuthResponsePayload, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}