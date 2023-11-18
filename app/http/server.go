package http

import (
	"context"
	"net/http"
	"time"
)

type server struct {
	server *http.Server
}

func New(handler http.Handler, port string) *server {
	server := &server{
		server: &http.Server{
			Addr:              ":" + port,
			Handler:           handler,
			ReadHeaderTimeout: 60 * time.Second,
		},
	}

	return server
}

func (s *server) ListenAndServe() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return
		}
	}()
}

func (s *server) Shutdown() {

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		return
	}
}