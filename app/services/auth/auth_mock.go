package auth

import (
	"context"

	"github.com/victor-felix/chat-service/app/models"
)

type ServiceMock struct {
	LoginFn func(ctx context.Context, payload *models.AuthPayload) (*models.AuthResponsePayload, error)
	LoginFnInvokedCount int
}

func (s *ServiceMock) Login(ctx context.Context, payload *models.AuthPayload) (*models.AuthResponsePayload, error) {
	s.LoginFnInvokedCount++
	return s.LoginFn(ctx, payload)
}