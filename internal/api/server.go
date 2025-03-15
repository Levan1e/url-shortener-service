package api

import (
	"context"
	"fmt"
	"net/http"

	http_helpers "github.com/Levan1e/url-shortener-service/pkg/http"
	"github.com/Levan1e/url-shortener-service/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type Handler interface {
	RegisterRoutes(r chi.Router)
}

type Server struct {
	server *http.Server
}

func NewServer(config *http_helpers.Config, handlers ...Handler) *Server {
	mux := chi.NewMux()
	mux.Route("/api",func(r chi.Router) {
		for _, handler := range handlers {
			handler.RegisterRoutes(r)
		}
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: mux,
	}
	return &Server{server: server}
}

func (s *Server) Run() error {
	logger.Infof("Start listen on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(_ context.Context) error {
	logger.Info("Got graceful shutdown signal, stopping server")
	return s.server.Close()
}
