package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"eridiumdev/yandex-praktikum-go-devops/config"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
)

type Server struct {
	Server *http.Server
}

func NewServer(handler http.Handler, cfg *config.ServerConfig) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    cfg.Address,
			Handler: handler,
			// golangci-lint: Potential Slowloris Attack because ReadHeaderTimeout is not configured in the http.Server
			ReadHeaderTimeout: time.Second,
		},
	}
}

func (s *Server) Start(ctx context.Context) {
	err := s.Server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		// ErrServerClosed is normal-case scenario (i.e. graceful server stop)
		logger.New(ctx).Fatalf("Failed to start HTTP server: %s", err.Error())
	}
}

func (s *Server) Stop(ctx context.Context) {
	err := s.Server.Shutdown(ctx)
	if err != nil {
		logger.New(ctx).Errorf("Error when stopping HTTP server: %s", err.Error())
	}
}
