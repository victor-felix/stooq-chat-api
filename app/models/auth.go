package models

import (
	"context"
	"strings"

	"github.com/victor-felix/chat-service/app/errors"
)

type AuthPayload struct {
	Email	string `json:"email"`
	Password	string `json:"password"`
}

type AuthResponsePayload struct {
	Token	string `json:"token"`
}

func (ap *AuthPayload) Validate() error {
	var field string

	field = strings.TrimSpace(ap.Email)
	if len(field) == 0 {
		return errors.NewErrInvalidArgument(ErrorRequiredFieldMissing).WithMessage("email is required")
	}

	field = strings.TrimSpace(ap.Password)
	if len(field) == 0 {
		return errors.NewErrInvalidArgument(ErrorRequiredFieldMissing).WithMessage("password is required")
	}

	return nil
}

type AuthService interface {
	Login(ctx context.Context, payload *AuthPayload) (*AuthResponsePayload, error)
}