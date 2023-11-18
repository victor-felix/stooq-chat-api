package brokers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/victor-felix/chat-service/app/models"
)

type StooqBroker struct {
	Receiver amqp.Queue
	Publisher amqp.Queue
	Channel *amqp.Channel
	Pool *models.Pool
	MessageService models.MessageService
	log zerolog.Logger
}

func (sb *StooqBroker) SetUp(
	receiverQueue,
	publisherQueue string,
	channel *amqp.Channel,
	wsPool *models.Pool,
	messageService models.MessageService,
	log zerolog.Logger,
) {
	queueReceiver, err := channel.QueueDeclare(
		receiverQueue,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	queuePublisher, err := channel.QueueDeclare(
		publisherQueue,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	sb.Receiver = queueReceiver
	sb.Publisher = queuePublisher
	sb.Channel = channel
	sb.Pool = wsPool
	sb.MessageService = messageService
}

func (sb *StooqBroker) Publish(requestBody chan []byte) {
	for body := range requestBody {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := sb.Channel.PublishWithContext(ctx, "", sb.Publisher.Name, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body: body,
		})
		cancel()
		if err != nil {
			sb.log.Error().Msg(err.Error())
			return
		}
	}
}

func (sb *StooqBroker) ReadMessages() {
	messages, err := sb.Channel.Consume(
		sb.Receiver.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		sb.log.Fatal().Msg(err.Error())
	}

	receivedMessages := make(chan models.Message)
	go sb.messageTransformer(messages, receivedMessages, sb.log)
	go sb.processResponse(receivedMessages, sb.Pool, sb.log)
}


func (sb *StooqBroker)  messageTransformer(entries <-chan amqp.Delivery, receivedMessages chan models.Message, log zerolog.Logger) {
	var botResponse models.Message
	for message := range entries {
		log.Info().Msg(fmt.Sprintf("Received a message: %s", string(message.Body)))

		err := json.Unmarshal([]byte(message.Body), &botResponse)

		if err != nil {
			log.Error().Msg(err.Error())
			continue
		}
		receivedMessages <- botResponse
	}
}

func (sb *StooqBroker) processResponse(responses <-chan models.Message, pool *models.Pool, log zerolog.Logger) {
	for response := range responses {
		log.Info().Msg(fmt.Sprintf("Processing bot response for room: %s, content: %s", response.RoomID, response.Content))

		response := models.Message{
			RoomID:  response.RoomID,
			Content: response.Content,
			UserID: response.UserID,
			CreatedAt: response.CreatedAt,
		}

		messageBody := models.Message{Content: response.Content, RoomID: response.RoomID, UserID: "", CreatedAt: response.CreatedAt}
		sb.MessageService.SaveMessageByAmqp(messageBody)
		message  := models.MessageWebsocket{Type: 1, Body: messageBody}

		pool.Broadcast <- message
	}
}
