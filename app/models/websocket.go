package models

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan MessageWebsocket
}

type Client struct {
	ID     string
	Pool   *Pool
	Conn   *websocket.Conn
	Email string
	UserID string
}

type MessageWebsocket struct {
	Type int `json:"type,omitempty"`
	Body Message `json:"body,omitempty"`
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan MessageWebsocket),
	}
}

func (p *Pool) Start() {
	defer p.RestartWebsocket()
	for {
		select {
		case client := <-p.Register:
			p.Clients[client] = true
			for client := range p.Clients {
				msgBody := Message{Content: "new user joined..."}
				client.Conn.WriteJSON(MessageWebsocket{Type: 1, Body: msgBody })
			}

		case client := <-p.Unregister:
			delete(p.Clients, client)
			for client := range p.Clients {
				msgBody := Message{
					Content: "user disconnected...",
				}
				client.Conn.WriteJSON(MessageWebsocket{Type: 1, Body: msgBody })
			}

		case message := <-p.Broadcast:
			for client := range p.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					return
				}
			}
		}
	}
}

func (p *Pool) RestartWebsocket() {
	if err := recover(); err != nil {
		go p.Start()
	}
}

type WebsocketService interface {
	Read(ctx context.Context, bodyChan chan []byte, client *Client)
	Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error)
}
