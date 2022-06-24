package routers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

type ChiRouter struct {
	Mux *chi.Mux
}

func NewChiRouter() *ChiRouter {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	return &ChiRouter{
		Mux: r,
	}
}

func (r *ChiRouter) GetHandler() http.Handler {
	return r.Mux
}

func (r *ChiRouter) AddRoute(method, endpoint string, handler http.HandlerFunc) {
	r.Mux.Method(method, endpoint, handler)
}

func (r *ChiRouter) URLParam(req *http.Request, name string) string {
	return chi.URLParam(req, name)
}
