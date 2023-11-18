package http

import (
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	configs "github.com/victor-felix/chat-service/app/config"
	"github.com/victor-felix/chat-service/app/models"
)

type Builder interface {
	SetUserService(userService models.UserService)
	SetMessageService(messageService models.MessageService)
	SetRoomService(roomService models.RoomService)
	SetAuthService(authService models.AuthService)
	SetWebsocketService(websocketService models.WebsocketService)
	SetAmqpConn(amqpConn *amqp.Connection)
	SetAmqpChannel(amqpChannel *amqp.Channel)
	SetPool(pool *models.Pool)
	SetConfig(config configs.Config)
	SetDevMode(devMode bool)
	SetLog(log zerolog.Logger)
	GetHandler() http.Handler
}

func NewHandlerBuilder() *handler {
	return &handler{}
}

func (h *handler) GetHandler() http.Handler {
	return NewHandler(
		h.userService,
		h.messageService,
		h.roomService,
		h.authService,
		h.websocketService,
		h.amqpConn,
		h.amqpChannel,
		h.pool,
		h.config,
		h.log,
		h.devMode,
	)
}

func (h *handler) SetUserService(userService models.UserService) {
	h.userService = userService
}

func (h *handler) SetMessageService(messageService models.MessageService) {
	h.messageService = messageService
}

func (h *handler) SetRoomService(roomService models.RoomService) {
	h.roomService = roomService
}

func (h *handler) SetAuthService(authService models.AuthService) {
	h.authService = authService
}

func (h *handler) SetWebsocketService(websocketService models.WebsocketService) {
	h.websocketService = websocketService
}

func (h *handler) SetAmqpConn(amqpConn *amqp.Connection) {
	h.amqpConn = amqpConn
}

func (h *handler) SetAmqpChannel(amqpChannel *amqp.Channel) {
	h.amqpChannel = amqpChannel
}

func (h *handler) SetPool(pool *models.Pool) {
	h.pool = pool
}

func (h *handler) SetConfig(config configs.Config) {
	h.config = config
}

func (h *handler) SetLog(log zerolog.Logger) {
	h.log = log
}

func (h *handler) SetDevMode(devMode bool) {
	h.devMode = devMode
}
