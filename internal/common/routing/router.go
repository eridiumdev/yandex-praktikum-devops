package routing

import (
	"net/http"
	"strings"
)

// Router bundles necessary routing-related functionality
type Router interface {
	GetHandler() http.Handler
	AddRoute(method, endpoint string, handler http.HandlerFunc)
	URLParam(req *http.Request, name string) string
}

func URLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimRight(r.URL.Path, "/")
		next.ServeHTTP(w, r)
	})
}
