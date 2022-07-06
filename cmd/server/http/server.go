package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"eridiumdev/yandex-praktikum-go-devops/cmd/server/http/routers"
	"eridiumdev/yandex-praktikum-go-devops/internal/logger"
)

type Server struct {
	Server *http.Server
}

func NewServer(router routers.Router, host string, port int, ) *Server {
	return &Server{
		Server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", host, port),
			Handler: router.GetHandler(),
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
