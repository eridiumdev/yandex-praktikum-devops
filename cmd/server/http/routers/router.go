package routers

import (
	"net/http"
)

type Router interface {
	GetHandler() http.Handler
	AddRoute(method, endpoint string, handler http.HandlerFunc)
	URLParam(req *http.Request, name string) string
}
