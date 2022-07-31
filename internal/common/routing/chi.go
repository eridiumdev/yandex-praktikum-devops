package routing

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

type ChiRouter struct {
	Mux *chi.Mux
}

func NewChiRouter(middlewares ...func(http.Handler) http.Handler) *ChiRouter {
	r := chi.NewRouter()
	for _, m := range middlewares {
		r.Use(m)
	}
	return &ChiRouter{
		Mux: r,
	}
}

func (r *ChiRouter) GetHandler() http.Handler {
	return r.Mux
}

func (r *ChiRouter) AddRoute(method, endpoint string, handler http.HandlerFunc) {
	if endpoint != "/" {
		endpoint = strings.TrimRight(endpoint, "/")
	}
	r.Mux.Method(method, endpoint, handler)
}

func (r *ChiRouter) URLParam(req *http.Request, name string) string {
	return chi.URLParam(req, name)
}
