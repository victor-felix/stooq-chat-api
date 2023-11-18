package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"

	"github.com/victor-felix/chat-service/app/brokers"
	configs "github.com/victor-felix/chat-service/app/config"
	"github.com/victor-felix/chat-service/app/errors"
	"github.com/victor-felix/chat-service/app/models"
)

const ConnectionPoolAttribute = "wsConnectionPool"

type handler struct {
	userService models.UserService
	messageService models.MessageService
	roomService models.RoomService
	authService models.AuthService
	websocketService models.WebsocketService
	amqpConn *amqp.Connection
	amqpChannel *amqp.Channel
	pool *models.Pool
	config configs.Config
	log zerolog.Logger
	devMode bool
}

func NewHandler(
	userService models.UserService,
	messageService models.MessageService,
	roomService models.RoomService,
	authService models.AuthService,
	websocketService models.WebsocketService,
	amqpConn *amqp.Connection,
	amqpChannel *amqp.Channel,
	pool *models.Pool,
	config configs.Config,
	log zerolog.Logger,
	devMode bool,
) http.Handler {
	handler := &handler{
		userService: userService,
		messageService: messageService,
		roomService: roomService,
		authService: authService,
		websocketService: websocketService,
		amqpConn: amqpConn,
		amqpChannel: amqpChannel,
		pool: pool,
		config: config,
		log: log,
		devMode: devMode,
	}

	go pool.Start()
	var broker brokers.StooqBroker
	broker.SetUp(config.RabbitMQ.StooqReceiverQueue, config.RabbitMQ.StooqPublisherQueue, handler.amqpChannel, pool, messageService, handler.log)

	router := gin.Default()

	if handler.devMode {
		router.Use(DevCORSMiddleware())
	}

	users := router.Group("/users")
	{
		users.POST("", handler.postUser)
	}

	messages := router.Group("/messages")
	{
		messages.POST("", handler.authGuard, handler.postMessage)
	}

	rooms := router.Group("/rooms")
	{
		rooms.GET("", handler.authGuard, handler.getRooms)
		rooms.POST("", handler.authGuard, handler.postRoom)
		rooms.GET("/:id/messages", handler.authGuard, handler.getRoomMessages)
	}

	auth := router.Group("/auth")
	{
		auth.POST("", handler.postAuth)
	}

	ws := router.Group("/ws")
	{
		ws.GET("", handler.authGuard, func(c *gin.Context) {
			jwtClaims, ok := c.Get(JwtClaimsAttribute)

			if !ok {
				handler.responseProblem(c, errors.NewErrUnauthorized(models.ErrorAuthTokenMissing).WithMessage("auth token is missing"))
				c.Abort()
				return
			}

			wsconn, err := handler.websocketService.Upgrade(c.Writer, c.Request)
			if err != nil {
				handler.log.Error().Err(err).Msg("error upgrading websocket connection")
				handler.responseProblem(c, err)
				c.Abort()
				return
			}

			client := &models.Client{
				Conn: wsconn,
				Pool: pool,
				Email: jwtClaims.(jwt.MapClaims)["email"].(string),
				UserID: jwtClaims.(jwt.MapClaims)["user_id"].(string),
				ID: jwtClaims.(jwt.MapClaims)["user_id"].(string),
			}

			pool.Register <- client
			requestBody := make(chan []byte)
			go handler.websocketService.Read(c.Request.Context(), requestBody, client)
			go broker.ReadMessages()
			go broker.Publish(requestBody)
		})
	}

	return router
}

func DevCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
