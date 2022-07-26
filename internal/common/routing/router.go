package routing

import "net/http"

// Router bundles necessary routing-related functionality
type Router interface {
	GetHandler() http.Handler
	AddRoute(method, endpoint string, handler http.HandlerFunc)
	URLParam(req *http.Request, name string) string
}
