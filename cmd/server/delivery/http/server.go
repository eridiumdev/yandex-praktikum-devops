package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"eridiumdev/yandex-praktikum-go-devops/internal/commons/logger"
)

// Router bundles necessary functionality for the Server to work
type Router interface {
	GetHandler() http.Handler
	AddRoute(method, endpoint string, handler http.HandlerFunc)
	URLParam(req *http.Request, name string) string
}

type Server struct {
	Server *http.Server
}

type ServerSettings struct {
	Host string
	Port int
}

func NewServer(router Router, settings ServerSettings) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", settings.Host, settings.Port),
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
