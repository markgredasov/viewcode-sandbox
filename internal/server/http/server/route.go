package server_http

import (
	"net/http"

	"github.com/markgredasov/viewcode-sandbox/internal/server/http/middleware"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []middleware.Middleware
}

func NewRoute(
	method,
	path string,
	handler http.HandlerFunc,
	middleware ...middleware.Middleware,
) *Route {
	return &Route{
		Method:     method,
		Path:       path,
		Handler:    handler,
		Middleware: middleware,
	}
}
