package routing

import (
	"net/http"
)

// Router bundles necessary routing-related functionality
type Router interface {
	GetHandler() http.Handler
	AddRoute(method, endpoint string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler)
	URLParam(req *http.Request, name string) string
}

func NotFound404(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("404 page not found"))
}
