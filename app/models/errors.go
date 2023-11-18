package models

import "github.com/victor-felix/chat-service/app/errors"

const (
	ErrorUnknown errors.ErrorCode = "error.unknown"
	ErrorRequiredFieldMissing errors.ErrorCode = "REQUIRED_FIELD_MISSING"
	ErrorEmailDuplicated errors.ErrorCode = "EMAIL_DUPLICATED"
	ErrorUserNotFound errors.ErrorCode = "USER_NOT_FOUND"
	ErrorInvalidPassword errors.ErrorCode = "INVALID_PASSWORD"
	ErrorInvalidToken errors.ErrorCode = "INVALID_TOKEN"
	ErrorAuthTokenMissing errors.ErrorCode = "AUTH_TOKEN_MISSING"
	ErrorRoomNotFound errors.ErrorCode = "ROOM_NOT_FOUND"
)

type ErrorDescriber interface {
	GetArgs() map[string]interface{}
	GetCode() errors.ErrorCode
	GetMessage() string
	Error() string
}

type ErrRulesNotSatisfied interface {
	ErrorDescriber
	ErrRulesNotSatisfied()
}

type ErrNotFound interface {
	ErrorDescriber
	ErrNotFound()
}

type ErrInvalidArgument interface {
	ErrorDescriber
	ErrInvalidArgument()
}

type ErrBadRequest interface {
	ErrorDescriber
	ErrBadRequest()
}

type ErrForbidden interface {
	ErrorDescriber
	ErrForbidden()
}

type ErrUnauthorized interface {
	ErrorDescriber
	ErrUnauthorized()
}

type ErrInternalServer interface {
	ErrorDescriber
	ErrInternalServer()
}

func ErrorDescriberCast(err error) (ErrorDescriber, bool) {
	if err == nil {
		return nil, false
	}
	e, ok := err.(ErrorDescriber)
	return e, ok
}

func ErrNotFoundCast(err error) (ErrNotFound, bool) {
	if err == nil {
		return nil, false
	}
	e, ok := err.(ErrNotFound)
	return e, ok
}

func ErrInvalidArgumentCast(err error) (ErrInvalidArgument, bool) {
	if err == nil {
		return nil, false
	}
	e, ok := err.(ErrInvalidArgument)
	return e, ok
}

func ErrRulesNotSatisfiedCast(err error) (ErrRulesNotSatisfied, bool) {
	if err == nil {
		return nil, false
	}
	e, ok := err.(ErrRulesNotSatisfied)
	return e, ok
}

func ErrBadRequestCast(err error) (ErrBadRequest, bool) {
	if err == nil {
		return nil, false
	}

	e, ok := err.(ErrBadRequest)
	return e, ok
}

func ErrForbiddenCast(err error) (ErrForbidden, bool) {
	if err == nil {
		return nil, false
	}

	e, ok := err.(ErrForbidden)
	return e, ok
}

func ErrUnauthorizedCast(err error) (ErrUnauthorized, bool) {
	if err == nil {
		return nil, false
	}

	e, ok := err.(ErrUnauthorized)
	return e, ok
}

func ErrInternalServerCast(err error) (ErrInternalServer, bool) {
	if err == nil {
		return nil, false
	}

	e, ok := err.(ErrInternalServer)
	return e, ok
}
