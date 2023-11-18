package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	config "github.com/victor-felix/chat-service/app/config"
	"github.com/victor-felix/chat-service/app/http"
	"github.com/victor-felix/chat-service/app/models"
	mongo "github.com/victor-felix/chat-service/app/pkg/mongo"
	"github.com/victor-felix/chat-service/app/services/auth"
	"github.com/victor-felix/chat-service/app/services/message"
	"github.com/victor-felix/chat-service/app/services/room"
	"github.com/victor-felix/chat-service/app/services/user"
	"github.com/victor-felix/chat-service/app/services/websocket"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	runLogFile, _ := os.OpenFile("app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	multi := zerolog.MultiLevelWriter(runLogFile)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	log.Info().Msg("starting application")
	
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Error().Msg(fmt.Sprintf("cannot load env variables %s", err))
	}

	conn, err := amqp.Dial(cfg.RabbitMQ.DSN)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("cannot connect to rabbitmq %q", err))
	}
	log.Info().Msg("connected to rabbitmq")
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Error().Msg(fmt.Sprintf("cannot create channel %q", err))
	}
	log.Info().Msg("channel created")
	defer channel.Close()
	
	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, cfg.MongoURL, log.Logger)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("cannot connect to mongo %q", err))
		return
	}
	log.Info().Msg("connected to mongo")

	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Error().Msg(fmt.Sprintf("cannot disconnect from mongo %q", err))
		}
	}()

	db := mongo.NewDatabase(mongoClient, cfg.DatabaseName)

	userStorage := mongo.NewUserStorage(db, log.Logger)
	messageStorage := mongo.NewMessageStorage(db, log.Logger)
	roomStorage := mongo.NewRoomStorage(db, log.Logger)

	userService := user.NewUserService(userStorage, cfg.JwtSecret, cfg.JwtTTL, log.Logger)
	messageService := message.NewMessageService(messageStorage)
	roomService := room.NewRoomService(roomStorage)
	authService := auth.NewAuthService(userService, cfg.JwtSecret, cfg.JwtTTL, log.Logger)
	websocketService := websocket.NewWebsocketService(messageService, log.Logger, cfg.Websocket.ReadBufferSize, cfg.Websocket.WriteBufferSize)

	handler := http.NewHandlerBuilder()
	handler.SetConfig(cfg)
	handler.SetUserService(userService)
	handler.SetMessageService(messageService)
	handler.SetRoomService(roomService)
	handler.SetAuthService(authService)
	handler.SetWebsocketService(websocketService)
	handler.SetDevMode(cfg.DevMode)
	handler.SetAmqpConn(conn)
	handler.SetAmqpChannel(channel)
	handler.SetPool(models.NewPool())

	server := http.New(handler.GetHandler(), cfg.Server.Port)

	server.ListenAndServe()

	log.Info().Msg("application started")

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan
	server.Shutdown()
}