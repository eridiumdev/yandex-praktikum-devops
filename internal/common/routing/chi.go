package routing

import (
	"net/http"

	"github.com/go-chi/chi"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/middleware"
)

type chiRouter struct {
	Mux *chi.Mux
}

func NewChiRouter(globalMiddlewares ...func(http.Handler) http.Handler) *chiRouter {
	r := chi.NewRouter()
	for _, m := range globalMiddlewares {
		r.Use(m)
	}

	// Add 404 not found handler, with basic logging middleware
	r.With(middleware.BasicSet...).NotFound(NotFound404)

	return &chiRouter{
		Mux: r,
	}
}

func (r *chiRouter) GetHandler() http.Handler {
	return r.Mux
}

func (r *chiRouter) AddRoute(
	method, endpoint string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	r.Mux.With(middlewares...).Method(method, endpoint, handler)
}

func (r *chiRouter) URLParam(req *http.Request, name string) string {
	return chi.URLParam(req, name)
}
