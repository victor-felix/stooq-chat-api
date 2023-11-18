package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/victor-felix/chat-service/app/models"
)

type WebsocketService struct {
	messageService models.MessageService
	log zerolog.Logger
	readBufferSize int
	writeBufferSize int
}

func NewWebsocketService(messageService models.MessageService, log zerolog.Logger, readBufferSize, writeBufferSice int) models.WebsocketService {
	return &WebsocketService{
		messageService: messageService,
		log: log,
		readBufferSize: readBufferSize,
		writeBufferSize: writeBufferSice,
	}
}

func (ws *WebsocketService) Read(ctx context.Context, bodyChan chan []byte, client *models.Client) {
	defer func ()  {
		client.Pool.Unregister <- client
		client.Conn.Close()
	}()

	defer client.Pool.RestartWebsocket()

	for {
		messageType, p, err := client.Conn.ReadMessage()
		if err != nil {
			ws.log.Error().Msg("error: " + err.Error())
			return
		}

		var msgBody models.Message
		err = json.Unmarshal(p, &msgBody)
		if err != nil {
			ws.log.Error().Msg("error: " + err.Error())
			return
		}

		msgBody.UserID = client.UserID

		msg := models.MessageWebsocket{Type: messageType, Body: msgBody}
		client.Pool.Broadcast <- msg
		ws.log.Info().Msg(fmt.Sprintf("Message Received: %+v", msg))

		if strings.Index(msgBody.Content, "/stock=") == 0 {
			bodyChan <- p
		} else {
			ws.log.Info().Msg(fmt.Sprintf("Message to Save: %+v", msgBody))
			go ws.messageService.SaveMessageByWebSocket(msgBody.RoomID, msgBody.UserID, msgBody.Content)
		}
	}
}

func (ws *WebsocketService) Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize: ws.readBufferSize,
		WriteBufferSize: ws.writeBufferSize,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	return conn, err
}
