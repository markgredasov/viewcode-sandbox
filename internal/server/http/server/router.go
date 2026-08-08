package server_http

import (
	"fmt"
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/internal/server/http/middleware"
)

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("v1")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
}

func NewApiVersionRouter(apiVersion ApiVersion) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.Handle(pattern, middleware.ChainMiddleware(route.Handler, route.Middleware...))
	}
}
