package app

import (
	"github.com/markgredasov/viewcode-sandbox/internal/server/http/middleware"
	server_http "github.com/markgredasov/viewcode-sandbox/internal/server/http/server"
)

func (a *App) initHTTPServer(routes []server_http.Route) error {
	apiVersion1Router := server_http.NewApiVersionRouter(server_http.ApiVersion1)
	apiVersion1Router.RegisterRoutes(routes...)

	serverConfig := server_http.NewConfigMust()
	server := server_http.NewHTTPServer(
		serverConfig,
		a.logger,
		middleware.RequestID(),
		middleware.Logger(a.logger),
		middleware.Recovery(),
		middleware.Trace(),
	)
	server.RegisterAPIRouters(apiVersion1Router)
	a.server = server

	return nil
}
