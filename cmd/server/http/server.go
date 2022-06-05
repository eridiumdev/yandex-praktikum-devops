package http

import (
	"eridiumdev/yandex-praktikum-go-devops/internal/logger"
	"fmt"
	"net/http"
)

type Server struct {
	Server *http.Server
}

type Middleware func(http.ResponseWriter, *http.Request) (stop bool)

func NewServer(host string, port int) *Server {
	return &Server{
		Server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", host, port),
		},
	}
}

func (s *Server) AddHandler(endpoint, method string, handler http.HandlerFunc) {
	http.Handle(endpoint, s.addMiddlewares(handler, s.filterByMethod(method), s.logRequest()))
}

func (s *Server) Start() {
	logger.Fatalf("Failed to start HTTP server: %s", s.Server.ListenAndServe().Error())
}

func (s *Server) addMiddlewares(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stop := false
		for _, middleware := range middlewares {
			stop = middleware(w, r)
			if stop {
				break
			}
		}
		handler(w, r)
	}
}

func (s *Server) filterByMethod(method string) Middleware {
	return func(w http.ResponseWriter, r *http.Request) (stop bool) {
		if r.Method != method {
			stop = true
		}
		return
	}
}

func (s *Server) logRequest() Middleware {
	return func(w http.ResponseWriter, r *http.Request) (stop bool) {
		logger.Infof("[server] new request: %s %s", r.Method, r.URL)
		return
	}
}
