package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"eridiumdev/yandex-praktikum-go-devops/internal/commons/logger"
)

type Server struct {
	Server *http.Server
}

type ServerSettings struct {
	Host string
	Port int
}

func NewServer(handler http.Handler, settings ServerSettings) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", settings.Host, settings.Port),
			Handler: handler,
		},
	}
}

func (s *Server) Start() {
	err := s.Server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		// ErrServerClosed is normal-case scenario (i.e. graceful server stop)
		logger.Fatalf("Failed to start HTTP server: %s", err.Error())
	}
}

func (s *Server) Stop(ctx context.Context) {
	err := s.Server.Shutdown(ctx)
	if err != nil {
		logger.Errorf("Error when stopping HTTP server: %s", err.Error())
	}
}
